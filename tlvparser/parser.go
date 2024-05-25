package tlvparser

/*
Tag - 8 bits/ 1 byte 0-256
Length - 8 bits/ 1 byte 0-256
Value - 256 bytes

Only Basic types to be supported
*/

func Encode(v interface{}) ([]byte, error) {
	var encoder TLVEncoder
	return encoder.Marshal(v)
}

func Decode(b []byte, v interface{}) error {
	var decoder TLVDecoder
	return decoder.Unmarshal(b, v)
}
