package io

import (
	"io/ioutil"
)

// Serializable interface needs to be implemented if you
// want to enable disk io for the object
// https://golang.org/pkg/encoding/gob/
type Serializable interface {
	// Any object that wants be writable needs to implement this function
	MarshalBinary() ([]byte, error)
}

// Deserializable interface needs to be implemented if you
// want to enable disk io for the object
// https://golang.org/pkg/encoding/gob/
type Deserializable interface {
	// Any object that wants to be read from bytes must implement this function
	UnmarshalBinary([]byte) error
}

// WriteToFile on input serializable and filename, writes the
// object into the file
func WriteToFile(s Serializable, fname string) {
	data, err := s.MarshalBinary()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fname, data, 0777)
	if err != nil {
		panic(err)
	}
}

// ReadFromFile on input a deserializable object, returns a serializable object
func ReadFromFile(d Deserializable, fname string) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	err = d.UnmarshalBinary(bytes)
	if err != nil {
		panic(err)
	}
}
