package util_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/adithyabhatkajake/libchatter/util"
)

func TestBufferedIO(t *testing.T) {
	var testData [100]byte
	var outBuf [100]byte
	var b bytes.Buffer

	for i := 0; i < len(testData); i++ {
		testData[i] = byte(i)
	}

	br := bufio.NewReader(&b)
	bw := bufio.NewWriter(&b)

	err := util.BufferedWrite(bw, testData[:])
	if err != nil {
		t.Error("Error writing out")
	}

	length, err := util.BufferedRead(br, outBuf[:])
	if err != nil {
		t.Error("Error reading in.")
	}
	if length != uint64(len(testData)) {
		t.Error("Incorrect number of bytes read")
	}
	if testData != outBuf {
		t.Error("Buffers do not match, read fail")
		t.Error(testData)
		t.Error(outBuf)
	}
}
