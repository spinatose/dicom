package net

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/suyashkumar/dicom/pkg/uid"
)

type serviceUserStatus int

const (
	serviceUserInitial = iota
	serviceUserAssociationActive
	serviceUserClosed
)

type ServiceUser struct {
	label    string
	upcallCh chan upcallEvent

	lock *sync.Mutex
	cond *sync.Cond
	disp *serviceDispatcher

	status serviceUserStatus
	cm     interface{}
}

type ServiceUserOptions struct {
	CalledAETitle    string
	CallingAETitle   string
	SOPClasses       []string
	TransferSyntaxes []string
}

type ServiceUserOption func(*ServiceUserOptions)

func DefaultServiceUserOptions(opts ...ServiceUserOption) *ServiceUserOptions {
	options := &ServiceUserOptions{
		CalledAETitle:    "unknown",
		CallingAETitle:   "unknown",
		TransferSyntaxes: uid.StandardTransferSyntaxes,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

func (o *ServiceUserOptions) Validate() error {
	if len(o.CalledAETitle) == 0 {
		return errors.New("empty called aetitle")
	}

	if len(o.CallingAETitle) == 0 {
		return errors.New("empyt calling aetitle")
	}

	if len(o.SOPClasses) == 0 {
		return errors.New("empty sop classes")
	}

	if len(o.TransferSyntaxes) == 0 {
		return errors.New("no transfer syntaxes specified")
	}

	for i, tuid := range o.TransferSyntaxes {
		cuid, err := uid.CanonicalTransferSyntaxUID(tuid)
		if err != nil {
			return errors.Wrap(err, "invalid transfer syntax")
		}

		o.TransferSyntaxes[i] = cuid
	}

	return nil
}

func NewServiceUser(opts ...ServiceUserOption) (*ServiceUser, error) {
	options := DefaultServiceUserOptions(opts...)

	if err := options.Validate(); err != nil {
		return nil, err
	}

	label := newUIDSeq("user")
	lock := &sync.Mutex{}
	su := &ServiceUser{
		label:    label,
		upcallCh: make(chan upcallEvent, 128),
		lock:     lock,
		cond:     sync.NewCond(lock),
		status:   serviceUserInitial,
	}

	go su.eventLoop()

	return su, nil
}

func (su *ServiceUser) eventLoop() {
	for evt := range su.upcallCh {
		if evt.typ == upcallEventHandshakeCompleted {
			su.lock.Lock()
			doassert(su.cm == nil)
			su.status = serviceUserAssociationActive
			su.cond.Broadcast()
			su.cm = evt.cm
			doassert(su.cm != nil)
			su.lock.Unlock()
		}

		doassert(evt.typ == upcallEventData)
		//su.disp.handleEvent(evt)
	}
}
