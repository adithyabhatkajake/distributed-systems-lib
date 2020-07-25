package log

import (
	"log"
)

var (
	debugLogger *log.Logger
)

// Logger instance that this library uses
type Logger struct{}

// Debug messages here
func Debug(s string) {
	debugLogger.Printf(s)
}

// Inspired from https://www.honeybadger.io/blog/golang-logging/
func init() {

}
