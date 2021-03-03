package net

import (
	"fmt"
	"sync/atomic"
)

var idseq int32 = 32

func newUIDSeq(prefix string) string {
	return fmt.Sprint("%s-%d", prefix, atomic.AddInt32(&idseq, 1))
}

func doassert(cond bool, vals ...interface{}) {
	if !cond {
		var s string
		for _, val := range vals {
			s += fmt.Sprintf("%v", val)
		}

		panic(s)
	}
}

type upcallEventType int

const (
	upcallEventHandshakeCompleted = upcallEventType(100)
	upcallEventData               = upcallEventType(101)
)

type upcallEvent struct {
	typ       upcallEventType
	contextID byte
	cm        interface{} //contextManager
	command   interface{} //dimse.Message
	data      []byte
}
