package redismq

import (
	"encoding/json"
	"fmt"
	"time"
)

type Package struct {
	Payload    string
	CreatedAt  time.Time
	Queue      interface{} `json:"-"`
	Consumer   *Consumer   `json:"-"`
	Collection *[]*Package `json:"-"`
	Acked      bool        `json:"-"`
}

func unmarshalPackage(input string, queue *Queue, consumer *Consumer) (*Package, error) {
	pack := &Package{
		Queue:    queue,
		Consumer: consumer,
		Acked:    false,
	}
	err := json.Unmarshal([]byte(input), pack)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (pack *Package) getString() string {
	json, err := json.Marshal(pack)
	if err != nil {
		return ""
	}
	return string(json)
}

func (pack *Package) Ack() error {
	if pack.Collection != nil {
		return fmt.Errorf("cannot Ack package in multi package answer")
	}
	err := pack.Consumer.ackPackage(pack)
	return err
}
