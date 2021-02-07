package beater

import (
	"fmt"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/iwanteggroll/compliancebeat/check"
	"github.com/iwanteggroll/compliancebeat/config"
	"github.com/iwanteggroll/compliancebeat/scripting"
)

// compliancebeat configuration.
type compliancebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of compliancebeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &compliancebeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts compliancebeat.
func (bt *compliancebeat) Run(b *beat.Beat) error {
	logp.Info("compliancebeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	for _, checkConfig := range bt.config.Checks {
		// take this block out
		// fmt.Println(checkConfig.Name)
		// fmt.Println(checkConfig.Path)
		// fmt.Println(checkConfig.Period)
		// fmt.Println(checkConfig.Category)
		// fmt.Println(checkConfig.Enabled)
		// fmt.Println(checkConfig.Params)
		// fmt.Println(len(checkConfig.Params))
		// take block out

		var sc scripting.Script
		checkInstance := check.ComplianceCheck{}
		sc = checkInstance.Setup(&checkConfig)

		go checkInstance.Run(func(events []beat.Event) {
			bt.client.PublishAll(events)
		}, sc)

	}

	for {
		select {
		case <-bt.done:
			return nil
		}
	}

	return nil

	// ticker := time.NewTicker(bt.config.Period)
	// counter := 1
	// for {
	// 	select {
	// 	case <-bt.done:
	// 		return nil
	// 	case <-ticker.C:
	// 	}

	// 	event := beat.Event{
	// 		Timestamp: time.Now(),
	// 		Fields: common.MapStr{
	// 			"type":    b.Info.Name,
	// 			"counter": counter,
	// 		},
	// 	}
	// 	bt.client.Publish(event)
	// 	logp.Info("Event sent")
	// 	counter++
	// }
}

// Stop stops compliancebeat.
func (bt *compliancebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
