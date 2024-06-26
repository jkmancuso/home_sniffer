package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type kafkaConfig struct {
	topic     string
	partition int
	transport string
	host      string
	port      int
}

type kafkaStore struct {
	cfg  kafkaConfig
	conn *kafka.Conn
}

func NewKafkaStore(ctx context.Context) (kafkaStore, error) {
	kafkaCfg := newKafkaCfg(ctx)

	kStore := kafkaStore{
		cfg: kafkaCfg,
	}

	conn, err := kafkaCfg.connectKafka(ctx)

	if err != nil {
		log.Errorf("Err: %v\ncould not connect to kafka with params: %+v", err, kafkaCfg)
		return kStore, err
	}

	kStore.setConn(conn)

	log.Debugf("Returning kafka store: %+v", kStore)

	return kStore, nil
}

func (store *kafkaStore) setConn(conn *kafka.Conn) {
	store.conn = conn
}

func newKafkaCfg(_ context.Context) kafkaConfig {
	loadEnv()

	topic := os.Getenv("KAFKA_TOPIC")
	partition, _ := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
	transport := os.Getenv("KAFKA_TRANSPORT")
	host := os.Getenv("KAFKA_HOST")
	port, _ := strconv.Atoi(os.Getenv("KAFKA_PORT"))

	kafkaCfg := kafkaConfig{
		topic:     topic,
		partition: partition,
		transport: transport,
		host:      host,
		port:      port,
	}

	log.Debugf("Loading kafka cfg: %+v\n", kafkaCfg)

	return kafkaCfg

}

// implement io.writer
func (store kafkaStore) Write(data []byte) (int, error) {
	var packets []packetData
	var kafkaMsgs []kafka.Message
	var err error
	var msgBytes []byte
	var numBytes int

	if err := json.Unmarshal(data, &packets); err != nil {
		log.Errorf("Error unmarshalling payload, %v", err)
		return 0, err
	}

	for _, payload := range packets {
		msgBytes, err = json.Marshal(payload)

		if err != nil {
			log.Errorf("Error marshalling payload, %v", err)
			return 0, err
		}

		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: msgBytes})
	}

	numBytes, err = store.conn.WriteMessages(kafkaMsgs...)

	if err != nil {
		log.Errorf("Error writing msg to kafka: %v", err)
		return 0, nil
	}

	log.Debugf("Successfully wrote %d bytes of kafka msg", numBytes)

	return numBytes, nil
}

// return a kafka connection handle
func (cfg *kafkaConfig) connectKafka(ctx context.Context) (*kafka.Conn, error) {

	log.Debugf("Connecting to kafka\n%+v", cfg)
	conn, err := kafka.DialLeader(ctx,
		cfg.transport,
		fmt.Sprintf("%s:%d", cfg.host, cfg.port),
		cfg.topic,
		cfg.partition)

	if err != nil {
		log.Errorf("Error connecting to kakfa: %v", err)
		return nil, err
	}

	log.Debugf("Success")

	return conn, err

}
