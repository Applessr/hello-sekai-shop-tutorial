package queue

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

func ConnectProducer(brokerUrl []string, apiKey, secret string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	if apiKey != "" && secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokerUrl, config)
	if err != nil {
		log.Printf("Error: Kafka producer connection failed: %s", err.Error())
		return nil, errors.New("error: Kafka producer connection failed")
	}
	return producer, nil
}

func PushMessageWithKeyToQueue(brokerUrl []string, apiKey, secret, topic, key string, message []byte) error {
	producer, err := ConnectProducer(brokerUrl, apiKey, secret)
	if err != nil {
		log.Printf("Error: Kafka producer connection failed: %s", err.Error())
		return errors.New("error: Kafka producer connection failed")
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(key),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error: Kafka producer failed to send message: %s", err.Error())
		return errors.New("error: Kafka producer failed to send message")
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

func ConnectConsumer(brokerUrl []string, apiKey, secret string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	if apiKey != "" && secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	consumer, err := sarama.NewConsumer(brokerUrl, config)
	if err != nil {
		log.Printf("Error: Kafka consumer connection failed: %s", err.Error())
		return nil, errors.New("error: Kafka consumer connection failed")
	}
	return consumer, nil
}

func DecodeMessage(obj any, value []byte) error {
	err := json.Unmarshal(value, obj)
	if err != nil {
		log.Printf("Error: Decode message failed: %s", err.Error())
		return errors.New("error: decode message failed")
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		log.Printf("Error: Validate message failed: %s", err.Error())
		return errors.New("error: validate message failed")
	}
	return nil
}
