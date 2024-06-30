package main

import (
	"context"
	"encoding/json"

	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/jkmancuso/home_sniffer/stores"
	pcap "github.com/packetcap/go-pcap"
	log "github.com/sirupsen/logrus"
)

type pcapConfig struct {
	device    string
	snaplen   int32
	promisc   bool
	timeout   time.Duration
	filter    string
	batchSize int
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
	filter := params["filter"]

	return pcapConfig{
		device:    device,
		snaplen:   int32(snaplen),
		promisc:   promisc,
		timeout:   time.Duration(timeout * int(time.Second)),
		filter:    filter,
		batchSize: batchSize,
	}
}

// Start new packet capture
func (cfg *pcapConfig) startPcap(store stores.Sender, cache *Cache, ctx context.Context) error {
	log.Debugf("Starting packet cap on device %v\n", cfg.device)

	handle, err := cfg.newPcapHandle()

	if err != nil {
		return err
	}

	if err = handle.SetBPFFilter(cfg.filter); err != nil {
		log.Fatalf("Unable to set filter %v", cfg.filter)
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

	var dnsPacket *layers.DNS
	var ipPacket *layers.IPv4
	var src, dst string
	var size uint16

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

			log.Debugf("Getting packet: %v", packet.String())

			if dnsLayer = packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
				dnsPacket = dnsLayer.(*layers.DNS)
				for _, answer := range dnsPacket.Answers {
					log.Infof("DNS: %+v", string(answer.Name))
				}
			}

			if ipLayer = packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ipPacket, _ = ipLayer.(*layers.IPv4)
				size = ipPacket.Length
				src = ipPacket.SrcIP.String()
				dst = ipPacket.DstIP.String()
				log.Infof("IP: %+v %+v %+v", size, src, dst)

			}

			//some packets have no payload such as ACKs, just move to next iteration
			if size == 0 || len(src) == 0 || len(dst) == 0 {
				//log.Debug("Dropping")
				continue
			}

			i += 1

			srcEntry, _ = NewIPinfo(src, *cache, ctx)
			dstEntry, _ = NewIPinfo(dst, *cache, ctx)
			entry, _ = NewEntryData(srcEntry, dstEntry, size)

			entryBytes, err = json.Marshal(entry)

			if err != nil {
				log.Errorf("Error marshalling %v\n%v", entry, err)
			}

			log.Debugf("Queueing up packet: %v", string(entryBytes))

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
func (cfg *pcapConfig) newPcapHandle() (*pcap.Handle, error) {

	if err := cfg.validateInterfaces(); err != nil {
		return nil, err
	}

	handle, err := pcap.OpenLive(
		cfg.device,
		cfg.snaplen,
		true,
		0, false)

	if err != nil {
		return nil, err
	}

	return handle, err

}

// Helper function to check that the selected device (ie br0 or eth0) is valid on your machine
func (cfg *pcapConfig) validateInterfaces() error {
	log.Println("Confirming valid devices")
	/*
		selectedDevice := cfg.device

		ifs, err := pcap.

		if err != nil {
			return err
		}

		for _, ethInterface := range ifs {
			if ethInterface.Name == selectedDevice {
				log.Printf("Found device %v, checking IP...\n", selectedDevice)

				ips := ethInterface.Addresses

				if len(ips) > 0 {
					log.Printf("IP: %v", ips)
					return nil
				}

			}
		}

		return errors.New("interface is not valid ")*/

	return nil

}

func NewEntryData(src ipInfo, dest ipInfo, size uint16) (entryData, error) {

	return entryData{
		Src:    src,
		Dst:    dest,
		Length: int(size),
	}, nil
}
