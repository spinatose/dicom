package dicom

import (
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/suyashkumar/dicom/pkg/uid"
)

type Options struct {
	ParseDataset          bool
	ParsePixelData        bool
	DefaultTransferSyntax string
	DefaultEncoding       string
	Limit                 int64
	IncludeTags           []tag.Tag
	FrameChannel          chan *frame.Frame
}

type Option func(*Options)

func DefaultOptions() *Options {
	return &Options{
		ParseDataset:          true,
		ParsePixelData:        true,
		DefaultTransferSyntax: uid.ExplicitVRLittleEndian,
	}
}

func NewOptions(opts ...Option) *Options {
	options := DefaultOptions()

	for _, o := range opts {
		o(options)
	}

	return options
}

func ParsePixelData(b bool) Option {
	return func(o *Options) {
		o.ParsePixelData = b
	}
}

func ParseDataset(b bool) Option {
	return func(o *Options) {
		o.ParseDataset = b
	}
}

func Limit(l int64) Option {
	return func(o *Options) {
		o.Limit = l
	}
}

func FrameChannel(c chan *frame.Frame) Option {
	return func(o *Options) {
		o.FrameChannel = c
	}
}

func IncludeTags(tags ...tag.Tag) Option {
	return func(o *Options) {
		o.IncludeTags = tags
	}
}
