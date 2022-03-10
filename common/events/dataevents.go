package events

import (
	"sync"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel es un canal que acepta un DataEvent
type DataChannel chan DataEvent

// DataChannelSlice es un slice de DataChannels
type DataChannelSlice []DataChannel

// EventBus almacena la informacion sobre suscriptores por topico
type EventBus struct {
	Subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.Subscribers[topic]; found {
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.Subscribers[topic]; found {
		eb.Subscribers[topic] = append(prev, ch)
	} else {
		eb.Subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}
