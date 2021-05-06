package dimse

import (
	"fmt"
	"io"
	"strconv"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/dicomio"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type MessageID uint16

type Message interface {
	Encode(*dicomio.Writer)
	MessageID() MessageID
	CommandField() int
	Status() *Status
	HasData() bool
}

type elementRequirement int

const (
	requiredElement elementRequirement = iota
	optionalElement
)

type messageDecoder struct {
	ds     dicom.Dataset
	parsed []bool
	err    error
}

func (d *messageDecoder) findElement(tag tag.Tag, req elementRequirement) *dicom.Element {
	for i, el := range d.ds.Elements {
		if el.Tag == tag {
			d.parsed[i] = true
			return el
		}
	}

	if req == requiredElement {
		d.err = fmt.Errorf("dimse.findElement: required element %s not found", tag.String())
	}

	return nil
}

func (d *messageDecoder) unparsedElements() []*dicom.Element {
	up := []*dicom.Element{}
	for i, parsed := range d.parsed {
		if !parsed {
			up = append(up, d.ds.Elements[i])
		}
	}

	return up
}

func (d *messageDecoder) status() Status {
	return Status{
		Status:       StatusCode(d.uintValue(tag.Status, requiredElement)),
		ErrorComment: d.stringValue(tag.ErrorComment, optionalElement),
	}
}

func (d *messageDecoder) stringValue(t tag.Tag, req elementRequirement) string {
	e := d.findElement(t, req)
	if e == nil {
		return ""
	}

	return e.Value.String()
}

func (d *messageDecoder) uintValue(t tag.Tag, req elementRequirement) uint16 {
	sv := d.stringValue(t, req)
	v, err := strconv.ParseUint(sv, 10, 16)
	if err != nil {
		return 0
	}

	return uint16(v)
}

func encodeElements(w io.Writer, ds dicom.Dataset) {

}
