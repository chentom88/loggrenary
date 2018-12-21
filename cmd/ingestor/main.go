package main

import (
	"context"
	"log"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator/rlp/app"
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
	config, err := app.LoadConfig()
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

	streamConn := loggregator.NewEnvelopeStreamConnector(config.LoggrAddr, tlsConfig)
	defer streamConn.Close()

	rx := streamConn.Stream(context.Background(), &loggregator_v2.EgressBatchRequest{
		ShardId:   config.ShardID,
		Selectors: selectors,
	})

	// Set up connection to RabbitMQ
	rabbitConn, err := amqp.Dial(config.RabbitAddr)

	for {
		batch := rx()

		for _, env := range batch {
		}
	}
}
