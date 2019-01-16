package main

import (
	"context"
	"fmt"
	"log"
	"os"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/chentom88/loggrenary/internal/ingestor"
)

var selectors = []*loggregator_v2.Selector{
	{
		SourceId: "doppler",
		Message: &loggregator_v2.Selector_Counter{
			Counter: &loggregator_v2.CounterSelector{},
		},
	},
	{
		SourceId: "doppler",
		Message: &loggregator_v2.Selector_Gauge{
			Gauge: &loggregator_v2.GaugeSelector{},
		},
	},
	{
		SourceId: "metron",
		Message: &loggregator_v2.Selector_Counter{
			Counter: &loggregator_v2.CounterSelector{},
		},
	},
	{
		SourceId: "metron",
		Message: &loggregator_v2.Selector_Gauge{
			Gauge: &loggregator_v2.GaugeSelector{},
		},
	},
}

func main() {
	config, err := ingestor.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// Set up connection to Loggregator
	tlsConfig, err := loggregator.NewEgressTLSConfig(
		config.CACertPath,
		config.CertPath,
		config.KeyPath,
	)
	if err != nil {
		log.Fatal("Failed to get tls config", err)
	}

	loggr := log.New(os.Stderr, "[", log.LstdFlags)
	streamConn := loggregator.NewEnvelopeStreamConnector(
		config.LoggrAddr,
		tlsConfig,
		loggregator.WithEnvelopeStreamLogger(loggr),
	)

	rx := streamConn.Stream(context.Background(), &loggregator_v2.EgressBatchRequest{
		ShardId:   config.ShardID,
		Selectors: selectors,
	})

	fmt.Println("Started ingestor")
	for {
		batch := rx()

		for _, env := range batch {
			fmt.Println("%+v \n", env)
		}
	}
	// Set up connection to RabbitMQ
	// rabbitConn, err := amqp.Dial(config.RabbitAddr)

	// for {
	// 	batch := rx()

	// 	for _, env := range batch {
	// 	}
	// }
}
