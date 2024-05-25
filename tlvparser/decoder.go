package tlvparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type TLVDecoder struct {
	tags      map[uint8]bool
	byteOrder binary.ByteOrder
}

func (d *TLVDecoder) Init(byteOrder binary.ByteOrder) {
	d.tags = make(map[uint8]bool)
	d.byteOrder = byteOrder
}

func (d *TLVDecoder) DecodeValue(b []byte, v interface{}) error {
	value := reflect.ValueOf(v)
	value = reflect.Indirect(value)
	valueType := reflect.TypeOf(value.Interface())

	switch value.Kind() {
	case reflect.Int8:
		tmp := int64(int8(b[0]))
		value.SetInt(tmp)
	case reflect.Int16:
		tmp := int64(int16(d.byteOrder.Uint16(b)))
		value.SetInt(tmp)
	case reflect.Int32:
		tmp := int64(int32(d.byteOrder.Uint32(b)))
		value.SetInt(tmp)
	case reflect.Int64:
		tmp := int64(d.byteOrder.Uint64(b))
		value.SetInt(tmp)
	case reflect.Int:
		tmp := int64(d.byteOrder.Uint32(b))
		value.SetInt(tmp)
	case reflect.Uint8:
		tmp := uint64(b[0])
		value.SetUint(tmp)
	case reflect.Uint16:
		tmp := uint64(d.byteOrder.Uint16(b))
		value.SetUint(tmp)
	case reflect.Uint32:
		tmp := uint64(d.byteOrder.Uint32(b))
		value.SetUint(tmp)
	case reflect.Uint64:
		tmp := d.byteOrder.Uint64(b)
		value.SetUint(tmp)
	case reflect.Uint:
		tmp := uint64(d.byteOrder.Uint32(b))
		value.SetUint(tmp)
	case reflect.String:
		value.SetString(string(b))
	case reflect.Slice:
		if value.IsNil() {
			value.Set(reflect.MakeSlice(value.Type(), 0, 1))
		}
		if valueType.Elem().Kind() == reflect.Uint8 {
			value.SetBytes(b)
		} else {
			return errors.New("value type `Slice of " + valueType.String() + "` is not support decode")
		}
	}
	return nil
}

func (d *TLVDecoder) ParseTLV(b []byte) (fragments, error) {
	tlvFragment := make(fragments)
	buffer := bytes.NewBuffer(b)

	var tag uint8
	var length uint8
	for {
		if err := binary.Read(buffer, d.byteOrder, &tag); err != nil {
			fmt.Printf("Binary Read error: %v", err)
		}
		if err := binary.Read(buffer, d.byteOrder, &length); err != nil {
			fmt.Printf("Binary Read error: %v", err)
		}
		value := make([]byte, length)
		if _, err := buffer.Read(value); err != nil {
			return nil, err
		}
		tlvFragment.Add(tag, value)
		if buffer.Len() == 0 {
			break
		}
	}
	return tlvFragment, nil
}

func (d *TLVDecoder) Unmarshal(b []byte, v interface{}) error {
	d.Init(binary.BigEndian)
	value := reflect.ValueOf(v)
	if value.Type().Kind() != reflect.Ptr || value.Elem().Type().Kind() != reflect.Struct {
		return errors.New("passed value is not a pointer to struct")
	}
	value = reflect.Indirect(value)

	// Parse through all the fields
	tlvFragment, err := d.ParseTLV(b)
	if err != nil {
		return err
	}
	// Fill the values
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)

		tag, omit, _, err := GetTLVTagValue(value.Type().Field(i))
		if err != nil {
			return err
		}

		if omit || !tlvFragment.Exists(tag) {
			continue
		}

		if d.tags[tag] {
			return fmt.Errorf("tag \"%v\" has been used for multiple fields", tag)
		}
		d.tags[tag] = true

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		} else if fieldValue.Kind() == reflect.Slice && fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 1))
		}

		if fieldValue.Kind() != reflect.Ptr {
			fieldValue = fieldValue.Addr()
		}
		err = d.DecodeValue(tlvFragment.Get(tag), fieldValue.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}
