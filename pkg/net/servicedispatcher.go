package net

import "sync"

type serviceDispatcher struct {
	label      string
	downcallCh chan stateEvent

	lock sync.Mutex
}

type stateEvent struct{}
