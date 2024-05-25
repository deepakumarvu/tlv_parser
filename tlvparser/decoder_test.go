package tlvparser_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/deepakumarvu/tlv_parser/testutils"
	"github.com/deepakumarvu/tlv_parser/tlvparser"

	"github.com/franela/goblin"
	"github.com/maxatome/go-testdeep/td"
)

// go test -coverprofile=coverage ./... ; go tool cover -html=coverage

const DecodeTestCaseFile = "test/decodeTestCases.json"

type DecodeTestCases struct {
	Name  string
	Cases []DecodeTestCase
}

type DecodeTestCase struct {
	Description string
	Input       []byte
	Output      struct {
		Decoded interface{}
		IsError bool
	}
}

func TestTLVDecodingType1(t *testing.T) {
	var testFile []DecodeTestCases
	type testStruct struct {
		X  int    `tlv:"0,omitempty"`
		X1 int64  `tlv:"1,omitempty"`
		Y  string `tlv:"2,omitempty"`
		O  string `tlv:"-"`
		Z  []byte `tlv:"3,omitempty"`
	}
	err := testutils.ReadJsonToStruct(testutils.GetFilePathInWD(DecodeTestCaseFile), &testFile)
	if err != nil {
		log.Fatal("Error reading the test case file:", err)
		return
	}
	testCases := testFile[0]
	g := goblin.Goblin(t)
	for _, testCase := range testCases.Cases {
		var in, exp testStruct
		err := tlvparser.Decode(testCase.Input, &in)
		fmt.Printf("in: %v\n", in)
		fmt.Printf("err: %v\n", err)
		g.Describe(testCases.Name, func() {
			g.It(fmt.Sprintf("%s - Error Validation", testCase.Description), func() {
				g.Assert(err != nil).Eql(testCase.Output.IsError)
			})
			if !testCase.Output.IsError {
				g.It(fmt.Sprintf("%s - Ouput Validation", testCase.Description), func() {
					testutils.InterfaceToStruct(testCase.Output.Decoded, &exp)
					g.Assert(td.Cmp(t, in, exp)).Eql(true)
				})
			}
		})
	}
}

func TestTLVDecodingAllTypes(t *testing.T) {
	var testFile []DecodeTestCases
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
	err := testutils.ReadJsonToStruct(testutils.GetFilePathInWD(DecodeTestCaseFile), &testFile)
	if err != nil {
		log.Fatal("Error reading the test case file:", err)
		return
	}
	testCases := testFile[1]
	g := goblin.Goblin(t)
	for _, testCase := range testCases.Cases {
		var in, exp testStruct
		err := tlvparser.Decode(testCase.Input, &in)
		fmt.Printf("in: %v\n", in)
		fmt.Printf("err: %v\n", err)
		g.Describe(testCases.Name, func() {
			g.It(fmt.Sprintf("%s - Error Validation", testCase.Description), func() {
				g.Assert(err != nil).Eql(testCase.Output.IsError)
			})
			if !testCase.Output.IsError {
				g.It(fmt.Sprintf("%s - Ouput Validation", testCase.Description), func() {
					testutils.InterfaceToStruct(testCase.Output.Decoded, &exp)
					g.Assert(td.Cmp(t, in, exp)).Eql(true)
				})
			}
		})
	}
}
