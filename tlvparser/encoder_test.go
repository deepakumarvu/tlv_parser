package tlvparser_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/deepakumarvu/tlv_parser/testutils"
	"github.com/deepakumarvu/tlv_parser/tlvparser"

	"github.com/franela/goblin"
)

// go test -coverprofile=coverage ./... ; go tool cover -html=coverage

const EncodeTestCaseFile = "test/encodeTestCases.json"

type EncodeTestCases struct {
	Name  string
	Cases []EncodeTestCase
}

type EncodeTestCase struct {
	Description string
	Input       interface{}
	Output      struct {
		Encoded string
		IsError bool
	}
}

func TestTLVEncodingType1(t *testing.T) {
	var testFile []EncodeTestCases
	type test1 struct {
		X  int    `tlv:"0,omitempty"`
		X1 int64  `tlv:"1,omitempty"`
		Y  string `tlv:"2,omitempty"`
		O  string `tlv:"-"`
		Z  []byte `tlv:"3,omitempty"`
	}
	err := testutils.ReadJsonToStruct(testutils.GetFilePathInWD(EncodeTestCaseFile), &testFile)
	if err != nil {
		log.Fatal("Error reading the test case file:", err)
		return
	}
	testCases := testFile[0]
	g := goblin.Goblin(t)
	for _, testCase := range testCases.Cases {
		var in test1
		testutils.InterfaceToStruct(testCase.Input, &in)
		out, err := tlvparser.Encode(in)
		out64 := base64.StdEncoding.EncodeToString(out)
		fmt.Printf("out: %v\n", out)
		fmt.Printf("out64: %v\n", out64)
		fmt.Printf("err: %v\n", err)
		g.Describe(testCases.Name, func() {
			g.It(fmt.Sprintf("%s - Error Validation", testCase.Description), func() {
				g.Assert(err != nil).Eql(testCase.Output.IsError)
			})
			if !testCase.Output.IsError {
				g.It(fmt.Sprintf("%s - Ouput Validation", testCase.Description), func() {
					g.Assert(out64).Eql(testCase.Output.Encoded)
				})
			}
		})
	}
}

func TestTLVEncodingType2(t *testing.T) {
	type testStruct struct {
		X int `tlv:"omitempty"`
	}
	g := goblin.Goblin(t)
	var in testStruct
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid Struct tag", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType3(t *testing.T) {
	type testStruct struct {
		X int
	}
	g := goblin.Goblin(t)
	var in testStruct
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid Struct tag", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType4(t *testing.T) {
	type testStruct struct {
		X int `tlv:"-1,omitempty"`
	}
	g := goblin.Goblin(t)
	var in testStruct
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid Struct tag", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType5(t *testing.T) {
	type testStruct struct {
		X struct {
			y int
		} `tlv:"1"`
	}
	g := goblin.Goblin(t)
	var in testStruct
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid field type", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType6(t *testing.T) {
	g := goblin.Goblin(t)
	_, err := tlvparser.Encode(1)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid input type", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType7(t *testing.T) {
	type testStruct struct {
		X []int `tlv:"1"`
	}
	g := goblin.Goblin(t)
	var in testStruct
	in.X = []int{12, 121}
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Invalid field slice type", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType8(t *testing.T) {
	type testStruct struct {
		X int `tlv:"1"`
		Y int `tlv:"1"`
	}
	g := goblin.Goblin(t)
	var in testStruct
	_, err := tlvparser.Encode(in)
	fmt.Printf("err: %v\n", err)
	g.Describe("Confliting tags", func() {
		g.It("Error Validation", func() {
			g.Assert(err != nil).Eql(true)
		})
	})
}

func TestTLVEncodingType9(t *testing.T) {
	type testStruct struct {
		X     int     `tlv:"1,omitempty"`
		X8    int8    `tlv:"2,omitempty"`
		X16   int16   `tlv:"3,omitempty"`
		X32   int32   `tlv:"4,omitempty"`
		X64   int64   `tlv:"5,omitempty"`
		UX    uint    `tlv:"6,omitempty"`
		UX8   uint8   `tlv:"7,omitempty"`
		UX16  uint16  `tlv:"8,omitempty"`
		UX32  uint32  `tlv:"9,omitempty"`
		UX64  uint64  `tlv:"10,omitempty"`
		STR   string  `tlv:"11,omitempty"`
		BArr  []byte  `tlv:"12,omitempty"`
		PX    *int    `tlv:"13,omitempty"`
		PX8   *int8   `tlv:"14,omitempty"`
		PX16  *int16  `tlv:"15,omitempty"`
		PX32  *int32  `tlv:"16,omitempty"`
		PX64  *int64  `tlv:"17,omitempty"`
		PUX   *uint   `tlv:"18,omitempty"`
		PUX8  *uint8  `tlv:"19,omitempty"`
		PUX16 *uint16 `tlv:"20,omitempty"`
		PUX32 *uint32 `tlv:"21,omitempty"`
		PUX64 *uint64 `tlv:"22,omitempty"`
		PSTR  *string `tlv:"23,omitempty"`
		PBArr *[]byte `tlv:"24,omitempty"`
	}

	var testFile []EncodeTestCases
	err := testutils.ReadJsonToStruct(testutils.GetFilePathInWD(EncodeTestCaseFile), &testFile)
	if err != nil {
		log.Fatal("Error reading the test case file:", err)
		return
	}
	testCases := testFile[1]
	g := goblin.Goblin(t)
	for _, testCase := range testCases.Cases {
		var in testStruct
		testutils.InterfaceToStruct(testCase.Input, &in)
		fmt.Printf("in.PBArr: %v\n", in.PBArr)
		out, err := tlvparser.Encode(in)
		out64 := base64.StdEncoding.EncodeToString(out)
		fmt.Printf("out: %v\n", out)
		fmt.Printf("out64: %v\n", out64)
		fmt.Printf("err: %v\n", err)
		g.Describe(testCases.Name, func() {
			g.It(fmt.Sprintf("%s - Error Validation", testCase.Description), func() {
				g.Assert(err != nil).Eql(testCase.Output.IsError)
			})
			if !testCase.Output.IsError {
				g.It(fmt.Sprintf("%s - Ouput Validation", testCase.Description), func() {
					g.Assert(out64).Eql(testCase.Output.Encoded)
				})
			}
		})
	}
}
