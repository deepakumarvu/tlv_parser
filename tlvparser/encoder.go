package tlvparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
)

type TLVEncoder struct {
	encodedTLV *bytes.Buffer
	tags       map[uint8]bool
	byteOrder  binary.ByteOrder
}

func (e *TLVEncoder) Init(byteOrder binary.ByteOrder) {
	e.encodedTLV = new(bytes.Buffer)
	e.tags = make(map[uint8]bool)
	e.byteOrder = byteOrder
}

func (e *TLVEncoder) GetBytes() []byte {
	return e.encodedTLV.Bytes()
}

func (e *TLVEncoder) Build(tag uint8, value reflect.Value, omitempty bool) error {
	value = reflect.Indirect(value)
	if (value.IsZero() && omitempty) || !value.CanInterface() {
		return nil
	}
	if e.tags[tag] {
		return fmt.Errorf("tag \"%v\" has been used for multiple fields", tag)
	}
	e.tags[tag] = true
	switch value.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := e.AddBlock(tag, TypeToSizeMap[value.Kind()], value.Interface()); err != nil {
			return err
		}
	case reflect.Int, reflect.Uint:
		value = value.Convert(reflect.TypeOf(int32(1)))
		if err := e.AddBlock(tag, TypeToSizeMap[value.Kind()], value.Interface()); err != nil {
			return err
		}
	case reflect.String:
		str := value.String()
		if err := e.AddBlocks(tag, len(str), []byte(str)); err != nil {
			return err
		}
	case reflect.Slice:
		// []byte
		if value.Type().Elem().Kind() == reflect.Uint8 {
			bys := value.Bytes()
			if err := e.AddBlocks(tag, len(bys), []byte(bys)); err != nil {
				return err
			}
		} else {
			return errors.New("value type " + value.Type().String() + " is not support TLV encode")
		}
	default:
		return errors.New("value type " + value.Type().String() + " is not support TLV encode")
	}
	return nil
}

func (e *TLVEncoder) AddBlock(tag uint8, length uint8, val any) error {
	if err := e.WriteBlock(tag); err != nil {
		log.Printf("Tag write error: %+v", err)
		return err
	}
	if err := e.WriteBlock(length); err != nil {
		log.Printf("Len write error: %+v", err)
		return err
	}
	if err := e.WriteBlock(val); err != nil {
		log.Printf("Binary write error: %+v", err)
		return err
	}
	return nil
}

func (e *TLVEncoder) AddBlocks(tag uint8, length int, val []byte) error {
	blocksToWrite := math.Ceil(float64(length) / TLVBLockSize)
	var si, ei int
	len := TLVBLockSize
	for i := 0; i < int(blocksToWrite); i++ {
		si = ei
		ei += TLVBLockSize
		if ei > length {
			ei = length
			len = length - si
		}
		if err := e.AddBlock(tag, uint8(len), val[si:ei]); err != nil {
			return err
		}
	}

	return nil
}

func (e *TLVEncoder) WriteBlock(val any) error {
	return binary.Write(e.encodedTLV, e.byteOrder, val)
}

func (e *TLVEncoder) Marshal(v interface{}) ([]byte, error) {
	e.Init(binary.BigEndian)
	value := reflect.ValueOf(v)
	value = reflect.Indirect(value)

	if value.Type().Kind() != reflect.Struct {
		return nil, errors.New("passed value is not a struct")
	}

	// Parse through all the fields
	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Field(i)
		if !hasValue(field) {
			continue
		}
		tagVal, omit, omitempty, err := GetTLVTagValue(value.Type().Field(i))
		if err != nil {
			return nil, err
		}
		if omit {
			continue
		}
		err = e.Build(tagVal, field, omitempty)
		if err != nil {
			return nil, err
		}
	}
	return e.GetBytes(), nil
}
