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

	br := bufio.NewReader(&b)
	bw := bufio.NewWriter(&b)

	util.BufferedWrite(bw, testData[:])
	util.BufferedRead(br, outBuf[:])
	if testData != outBuf {
		t.Error("Buffers do not match, read fail")
	}
}
