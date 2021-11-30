package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	lowLevelCall()
}

func highLevelCall() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "test",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer r.Close()

	m, err := r.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

}

func lowLevelCall() {
	topic := "test"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)
	defer batch.Close()

	b := make([]byte, 10e3)
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}
}
