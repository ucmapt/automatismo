package fsm

import (
	"errors"
	"sync"
)

var ErrEventRejected = errors.New("Evento rechazado")

const (
	Default StateType = ""

	NoAction EventType = "NA"
)

type StateType string

type EventType string

type EventContext interface{}

type Action interface {
	Execute(evContext EventContext) EventType
}

type Events map[EventType]StateType

type State struct {
	Action Action
	Events Events
}

type States map[StateType]State

type StateMachine struct {
	Previous StateType
	Current  StateType
	States   States
	mutex    sync.Mutex
}

func (s *StateMachine) getNextState(event EventType) (StateType, error) {
	if state, ok := s.States[s.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return Default, ErrEventRejected
}

func (s *StateMachine) SendEvent(event EventType, evContext EventContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		nextState, err := s.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		state, ok := s.States[nextState]
		if !ok || state.Action == nil {
			// error de configuración
		}

		s.Previous = s.Current
		s.Current = nextState

		nextEvent := state.Action.Execute(evContext)
		if nextEvent == NoAction {
			return nil
		}
		event = nextEvent
	}
}
