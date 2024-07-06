package main

import (
	"context"
	"encoding/json"
	"net"

	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"

	pcap "github.com/packetcap/go-pcap"
	log "github.com/sirupsen/logrus"

	"github.com/jkmancuso/home_sniffer/stores"
)

type pcapConfig struct {
	device    string
	snaplen   int32
	promisc   bool
	timeout   time.Duration
	filter    string
	batchSize int
	syscalls  bool
}

type entryData struct {
	Src    ipInfo
	Dst    ipInfo
	Length int
}

func NewPcapCfg(params map[string]string) pcapConfig {

	device := params["device"]
	snaplen, _ := strconv.Atoi(params["snaplen"])
	timeout, _ := strconv.Atoi(params["timeout"])
	batchSize, _ := strconv.Atoi(params["batch_size"])
	promisc, _ := strconv.ParseBool(params["promisc"])
	syscalls, _ := strconv.ParseBool(params["syscalls"])

	filter := params["filter"]

	return pcapConfig{
		device:    device,
		snaplen:   int32(snaplen),
		promisc:   promisc,
		timeout:   time.Duration(timeout * int(time.Second)),
		filter:    filter,
		batchSize: batchSize,
		syscalls:  syscalls,
	}
}

// Start new packet capture
func (cfg *pcapConfig) startPcap(ctx context.Context, store stores.Sender, cache Cache, m *metrics) error {
	log.Printf("Starting packet cap on device %v\n", cfg.device)

	handle, err := cfg.NewPcapHandle()

	if err != nil {
		log.Errorf("Unable to get handle to packet capture %v", err)
		return err
	}

	if err = handle.SetBPFFilter(cfg.filter); err != nil {
		log.Errorf("Unable to set filter %v", cfg.filter)
		return err
	}

	defer handle.Close()
	defer store.Teardown()

	packetSource := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)

	var i int64
	var entry entryData
	var srcEntry, dstEntry ipInfo
	var entryBytes []byte
	var entryBatch []string

	var dnsLayer gopacket.Layer
	var ipLayer gopacket.Layer

	var question layers.DNSQuestion
	var answer layers.DNSResourceRecord

	var dnsPacket *layers.DNS
	var ipPacket *layers.IPv4
	var src, dst string
	var size uint16

	var parsedIP net.IP

	var DNSname, resolvedIPAddress string

	batchSize := cfg.batchSize

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for packet := range packetSource.Packets() {

		select {

		case <-sigCh:
			_ = store.Send(entryBatch)
			log.Printf("Caught interrupt, finishing remaining processing: %d packets\n", len(entryBatch))
			return nil
		default:

			if dnsLayer = packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
				dnsPacket = dnsLayer.(*layers.DNS)

				for _, question = range dnsPacket.Questions {
					DNSname = string(question.Name)
					log.Infof("DNS Q: %v", DNSname)
				}

				for _, answer = range dnsPacket.Answers {
					resolvedIPAddress = string(answer.String())

					// if it returns anything other than an IP, like a CNAME
					if parsedIP = net.ParseIP(resolvedIPAddress); parsedIP == nil {
						continue
					}

					log.Infof("DNS A: %+v", resolvedIPAddress)

					if err = cache.Set(ctx, resolvedIPAddress, DNSname); err != nil {
						log.Errorf("Error setting cache key %v %v", resolvedIPAddress, err)
					}
				}

				continue
				//if you get DNS packet no need to process the rest of the ip layer
			}

			if ipLayer = packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ipPacket, _ = ipLayer.(*layers.IPv4)
				size = ipPacket.Length
				src = ipPacket.SrcIP.String()
				dst = ipPacket.DstIP.String()
			}

			//some packets have no payload such as ACKs, just move to next iteration
			if size == 0 || len(src) == 0 || len(dst) == 0 {
				continue
			}

			i += 1

			srcEntry, _ = NewIPinfo(ctx, src, cache)
			dstEntry, _ = NewIPinfo(ctx, dst, cache)
			entry, _ = NewEntryData(srcEntry, dstEntry, size)

			srcEntry.updateCacheMetrics(m)
			dstEntry.updateCacheMetrics(m)

			log.Infof("IP: src %v (%v) dest %v (%v)", src, srcEntry.DNS, dst, dstEntry.DNS)

			entryBytes, err = json.Marshal(entry)

			if err != nil {
				log.Errorf("Error marshalling %v\n%v", entry, err)
			}

			entryBatch = append(entryBatch, string(entryBytes))

			if i%int64(batchSize) == 0 {
				log.Debug("Writing batch")

				if err := store.Send(entryBatch); err != nil {
					log.Errorf("Error writing to kafka: %v", err)
				}

				entryBatch = entryBatch[:0]
			}

		}

	}

	return nil

}

// Generate a new pcap handle given a config
func (cfg *pcapConfig) NewPcapHandle() (*pcap.Handle, error) {

	handle, err := pcap.OpenLive(
		cfg.device,
		cfg.snaplen,
		cfg.promisc,
		cfg.timeout,
		cfg.syscalls)

	if err != nil {
		return nil, err
	}

	return handle, err

}

func NewEntryData(src ipInfo, dest ipInfo, size uint16) (entryData, error) {

	return entryData{
		Src:    src,
		Dst:    dest,
		Length: int(size),
	}, nil
}
