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

func NewPcapCfg(params map[string]*string) pcapConfig {

	device := *params["device"]
	snaplen, _ := strconv.Atoi(*params["snaplen"])
	timeout, _ := strconv.Atoi(*params["timeout"])
	batchSize, _ := strconv.Atoi(*params["batch_size"])
	promisc, _ := strconv.ParseBool(*params["promisc"])
	syscalls, _ := strconv.ParseBool(*params["syscalls"])

	filter := *params["filter"]

	cfg := pcapConfig{
		device:    device,
		snaplen:   int32(snaplen),
		promisc:   promisc,
		timeout:   time.Duration(timeout * int(time.Second)),
		filter:    filter,
		batchSize: batchSize,
		syscalls:  syscalls,
	}

	log.Printf("Using cfg: %+v", cfg)
	return cfg
}

// Start new packet capture
func (cfg *pcapConfig) startPcap(ctx context.Context, store stores.Sender, cache Cache) error {
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

	var src, dst string
	var size uint16

	batchSize := cfg.batchSize

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for packet := range packetSource.Packets() {

		select {

		case <-sigCh:
			log.Printf("Caught interrupt, finishing remaining processing: %d packets\n", len(entryBatch))
			_ = store.Send(entryBatch)
			return nil
		default:

			if handleDNSLayer(ctx, packet, cache) {
				continue
			}

			size, src, dst = handleIPLayer(packet)

			//some packets have no payload such as ACKs, just move to next iteration
			if size == 0 || len(src) == 0 || len(dst) == 0 {
				continue
			}

			i += 1

			srcEntry, _ = NewIPinfo(ctx, src, cache)
			dstEntry, _ = NewIPinfo(ctx, dst, cache)
			entry, _ = NewEntryData(srcEntry, dstEntry, size)

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

func handleDNSLayer(ctx context.Context, p gopacket.Packet, c Cache) bool {
	if dnsLayer := p.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dnsPacket := dnsLayer.(*layers.DNS)

		var question layers.DNSQuestion
		var answer layers.DNSResourceRecord

		var DNSname, resolvedIPAddress string
		var parsedIP net.IP

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

			if err := c.Set(ctx, resolvedIPAddress, DNSname); err != nil {
				log.Errorf("Error setting cache key %v %v", resolvedIPAddress, err)
			}
		}

		return true
	}
	return false
}

func handleIPLayer(p gopacket.Packet) (uint16, string, string) {
	var ipLayer gopacket.Layer
	var ipPacket *layers.IPv4

	if ipLayer = p.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ipPacket, _ = ipLayer.(*layers.IPv4)
		return ipPacket.Length, ipPacket.SrcIP.String(), ipPacket.DstIP.String()
	}

	return 0, "", ""
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

	entry := entryData{
		Src:    src,
		Dst:    dest,
		Length: int(size),
	}

	return entry, nil
}
