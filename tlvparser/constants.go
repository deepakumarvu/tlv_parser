package tlvparser

import "reflect"

const (
	TLVBLockSize  = 256
	OmitEmpty     = "omitempty"
	TLV           = "tlv"
	Ignore        = "-"
	TagsSeperator = ","
)

var TypeToSizeMap = map[reflect.Kind]uint8{
	reflect.Int8:   1,
	reflect.Int16:  2,
	reflect.Int32:  4,
	reflect.Int64:  8,
	reflect.Uint8:  1,
	reflect.Uint16: 2,
	reflect.Uint32: 4,
	reflect.Uint64: 8,
}
