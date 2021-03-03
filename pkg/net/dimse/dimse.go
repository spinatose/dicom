package dimse

import (
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/suyashkumar/dicom"
)

type Message interface {
	Encode(*dicomio.Encoder)
	MessageID() MessageID
	CommandField() int
	Status() *Status
	HasData() bool
}

type Status struct {
	Status       StatusCode
	ErrorComment string
}

type messageDecoder struct {
	elems  []*dicom.Element
	parsed []bool
	err    error
}
