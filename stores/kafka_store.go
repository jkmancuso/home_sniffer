package stores

import (
	"context"
	"crypto/tls"

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
	tlsConfig *tls.Config
}

type kafkaStore struct {
	cfg  kafkaConfig
	conn *kafka.Conn
}

func NewKafkaStore(ctx context.Context, tlsConfigs ...*tls.Config) (kafkaStore, error) {
	kafkaCfg := newKafkaCfg(ctx)

	if len(tlsConfigs) != 0 {
		log.Info("tls enabled")
		kafkaCfg.setTLS(tlsConfigs[0])
	}

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

func (store kafkaStore) Teardown() {
	log.Printf("Tearing down kafka store")
}

func (store *kafkaStore) setConn(conn *kafka.Conn) {
	store.conn = conn
}

func newKafkaCfg(_ context.Context) kafkaConfig {

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

func (store kafkaStore) Send(data []string) error {
	var kafkaMsgs []kafka.Message
	var err error
	var numBytes int

	for _, payload := range data {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: []byte(payload)})
	}

	numBytes, err = store.conn.WriteMessages(kafkaMsgs...)

	if err != nil {
		log.Errorf("Error writing msg to kafka: %v", err)
		return nil
	}

	log.Debugf("Successfully wrote %d bytes of kafka msg", numBytes)

	return nil
}

func (cfg *kafkaConfig) setTLS(tlsCfg *tls.Config) {
	cfg.tlsConfig = tlsCfg
}

// return a kafka connection handle
func (cfg *kafkaConfig) connectKafka(ctx context.Context) (*kafka.Conn, error) {

	log.Debugf("Connecting to kafka\n%+v", cfg)

	dialer := &kafka.Dialer{
		Timeout:   0,
		DualStack: true,
		TLS:       cfg.tlsConfig,
	}

	conn, err := dialer.DialLeader(ctx,
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
