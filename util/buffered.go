package util

import (
	"bufio"
	"encoding/binary"
)

// BufferedRead reads the length first, and then reads that many bytes from the
// buffer
func BufferedRead(reader *bufio.Reader, msgBuf []byte) (uint64, error) {
	len := uint64(0)
	var lengthBuf [8]byte
	_, err := reader.Read(lengthBuf[:])
	length := binary.BigEndian.Uint64(lengthBuf[:])
	if err != nil {
		return len, err
	}
	for bytesRead := uint64(0); bytesRead < length; {
		readLen, err := reader.Read(msgBuf[bytesRead:])
		if err != nil {
			return len, err
		}
		bytesRead += uint64(readLen)
		len += uint64(readLen)
	}
	return len, nil
}

// BufferedWrite writes the length first, and then writes that many bytes from
// the buffer
func BufferedWrite(writer *bufio.Writer, data []byte) error {
	var lengthBuf [8]byte
	binary.BigEndian.PutUint64(lengthBuf[:], uint64(len(data)))
	writer.Write(lengthBuf[:])
	writer.Flush()
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
