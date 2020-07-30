package e2c

import (
	"fmt"
)

func (e *E2C) react(m []byte) {
	fmt.Println("Received a message of size", len(m))
}
