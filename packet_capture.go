package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	device  string
	snaplen int32
	promisc bool
	timeout time.Duration
}

type packetData struct {
	Src    ipInfo
	Dst    ipInfo
	Length int
}

func NewPcapCfg(device string) pcapConfig {
	return pcapConfig{
		device:  device,
		snaplen: 1500,
		promisc: true,
		timeout: 0,
	}
}

// Start new packet capture
func (cfg *pcapConfig) startPcap(store stores.Sender, cache *Cache, ctx context.Context) error {
	log.Debugf("Starting packet cap on device %v\n", cfg.device)

	handle, err := cfg.newPcapHandle()

	if err != nil {
		return err
	}

	_ = handle.SetBPFFilter("dns")

	defer handle.Close()
	defer store.Teardown()

	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)

	var netLayer, transportLayer, srcIP, dstIP, size string
	var i int64
	var pack packetData
	var packBytes []byte
	var packetBatch []string

	batchSize := 100

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for packet := range packetSource.Packets() {

		select {

		case <-sigCh:
			_ = store.Send(packetBatch)
			log.Printf("Caught interrupt, finishing remaining processing: %d packets\n", len(packetBatch))
			return nil
		default:

			netLayer = fmt.Sprintf("%+v", packet.NetworkLayer().LayerPayload())
			transportLayer = fmt.Sprintf("%v", packet.ApplicationLayer().LayerPayload())
			log.Debugf("transport %v\n\n", transportLayer)

			srcIP, dstIP = parseIPs(netLayer)
			size = parseSize(transportLayer)

			//some packets have no payload such as ACKs, just move to next iteration
			if len(size) == 0 || len(srcIP) == 0 || len(dstIP) == 0 {
				//log.Debug("Dropping")
				continue
			}

			log.Debugf("src:%v,dst:%v,size:%v\n", srcIP, dstIP, size)

			i += 1

			sizeInt, _ := strconv.Atoi(size)

			src, errSrc := GetIPLookupInfo(srcIP, *cache, ctx)
			dst, errDst := GetIPLookupInfo(dstIP, *cache, ctx)

			if errSrc != nil || errDst != nil {
				log.Warnf("Error looking up ip info: %v %v", errSrc, errDst)
			}

			pack = packetData{
				Src:    src,
				Dst:    dst,
				Length: sizeInt,
			}

			packBytes, err = json.Marshal(pack)

			if err != nil {
				log.Errorf("Error marshalling %v\n%v", pack, err)
			}

			log.Debugf("Queueing up packet: %v", string(packBytes))

			packetBatch = append(packetBatch, string(packBytes))

			if i%int64(batchSize) == 0 {
				log.Debug("Writing batch")

				if err := store.Send(packetBatch); err != nil {
					log.Errorf("Error writing to kafka: %v", err)
				}

				packetBatch = packetBatch[:0]
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
