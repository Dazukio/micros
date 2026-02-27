package shared_kafka

import (
	"sync"
)

type InMemory struct {
	Events   []*Task
	mu       *sync.Mutex
	notifier chan struct{}
}

func NewInMemory() *InMemory {
	return &InMemory{Events: make([]*Task, 0), mu: &sync.Mutex{}, notifier: make(chan struct{}, 1)}
}

func (i *InMemory) WriteEvent(task *Task) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	select {
	case i.notifier <- struct{}{}:
	default:
	}
	i.Events = append(i.Events, task)
	return nil
}

func (i *InMemory) ReadEvents() (*Task, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.Events) == 0 {
		return nil, nil
	}
	event := i.Events[0]
	i.Events = i.Events[1:]
	return event, nil
}

func (i *InMemory) GetNotifier() <-chan struct{} {
	return i.notifier
}
