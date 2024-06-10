package main

// default is send packets to kafka topic but I might do other options
type packetStore interface {
	sendSingle(packetData) error
	sendBatch([]packetData) error
}
