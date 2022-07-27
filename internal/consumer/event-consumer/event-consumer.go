package eventconsumer

import (
	"home/internal/events"
	"log"
	"time"
)

type Consumer struct {
	Fetcher   events.Fetcher
	Processor events.Processor
	BatchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		Fetcher:   fetcher,
		Processor: processor,
		BatchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.Fetcher.Fetch(c.BatchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s\n", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			return err
		}

	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Println("got new event: " + event.Text)

		if err := c.Processor.Process(event); err != nil {
			log.Println("can't handle event: " + err.Error())

			continue
		}
	}
	return nil
}
