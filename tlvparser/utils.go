package tlvparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetTLVTagValue(field reflect.StructField) (uint8, bool, bool, error) {
	var omit, omitempty bool
	var tagVal int
	var err error
	tlvTag, hasTLV := field.Tag.Lookup(TLV)
	if !hasTLV || len(tlvTag) == 0 {
		return 0, omit, omitempty, fmt.Errorf("field \"%s\" needs to have a `tlv` tag", field.Name)
	}
	tags := strings.Split(tlvTag, TagsSeperator)
	if tags[0] == Ignore {
		omit = true
	} else {
		tagVal, err = strconv.Atoi(tags[0])
		if err != nil || tagVal < 0 {
			return 0, omit, omitempty, fmt.Errorf("invalid `tlv` tag \"%s\" for field \"%s\". This needs to be a valid positive integer", tags[0], field.Name)
		}
	}
	if len(tags) > 1 && tags[1] == OmitEmpty {
		omitempty = true
	}
	return uint8(tagVal), omit, omitempty, nil
}

func canReference(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

func hasValue(value reflect.Value) bool {
	if canReference(value.Type()) {
		return !value.IsNil()
	} else {
		return value.IsValid()
	}
}
