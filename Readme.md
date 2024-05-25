# TLV Parser

## Why?

In golang there are parser for various data formats. But TLV is one such format for which we don't seem to have wide variety of parsers. This parser is designed keeping low memory devices in mind. Hence, the Tag and Length are kept as 1 Byte.


**Note:** This parser support basic types only. If you want to encode complex values encode them to byte array and then encode it as TLV.

## What?

##### Sample TLV Object:

| 1 Byte | 1 Byte | Length Bytes(Max 256) |
| :----: | :----: | :-------------------: |
|  Tag   | Length |         Value         |

##### Serial TLV Objects:

| 1 Byte | 1 Byte | Length Bytes(Max 256) | 1 Byte | 1 Byte | Length Bytes(Max 256) |
| :----: | :----: | :-------------------: | :----: | :----: | :-------------------: |
|  Tag   | Length |         Value         |  Tag   | Length |         Value         |

### Tag

Tag is a one byte field whose value can be in range 0-255.

### Length

Length is a one byte field whose value can be 1-255.

### Value

Value is a byte array whose length can be maximum of 255.

## Edge Cases

### Omitempty
Always use `omitempty` in case of `[]byte` or `string` as their default value is of length 0.

### Long Values
When the value you are trying to encode whose length is greater than 255 bytes then this package auto handles it by appending the remaining data in the adjacent TLV block. Below is a visual representation of how that value of a tag `1` will be encoded when it greater than 255 bytes.

|  Tag  | Length |   Value   |  Tag  | Length |   Value   |
| :---: | :----: | :-------: | :---: | :----: | :-------: |
|   1   |  255   | 255 Bytes |   1   |  120   | 120 Bytes |

## Supported Types

- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `string`
- `[]byte`

And pointers of all the above types as well like `*int`, `*[]byte`.

## TLV Struct Tag

```go
type IotData struct {
	X int    `tlv:"1"`
	Y string `tlv:"-"`           // Ignored
	Z string `tlv:"2,omitempty"` // Ignored if the value is default value i.e. 0 in case of int, empty string in case of string.
}
```

## Example

```go
package main

import (
	"fmt"

	"github.com/deepakumarvu/tlv_parser/tlvparser"
)

type IotData struct {
	X int    `tlv:"0"`
	Y string `tlv:"-"`
}

func main() {
	x := IotData{X: 0, Y: "NotEncoded"}
	// Encode the data
	encodedData, err := tlvparser.Encode(x)
	if err != nil {
		fmt.Printf("Error encoding the data: %v", err)
		return
	}

	// Decode the data
	var receivedData IotData
	err = tlvparser.Decode(encodedData, &receivedData)
	if err != nil {
		fmt.Printf("Error decoding the data: %v", err)
		return
	}
	fmt.Printf("Decoded data: %v", receivedData)
}
```