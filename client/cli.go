package main

import (
	"context"
	"fmt"
	p2p "github.com/libp2p/go-libp2p"
)

func main() {
	fmt.Println ("I am the client")
	ctx := context.Background()

	node, err := p2p.New(ctx)
	if err != nil {
		panic(err)
	}

	
	option p2p.Option()
}
