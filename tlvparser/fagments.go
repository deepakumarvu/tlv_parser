package tlvparser

type fragments map[uint8][]byte

func (f fragments) Add(tag uint8, buf []byte) {
	f[tag] = append(f[tag], buf...)
}

func (f fragments) Get(tag uint8) []byte {
	ret := f[tag]
	return ret
}

func (f fragments) Exists(tag uint8) bool {
	ret, t := f[tag]
	return t && len(ret) > 0
}
