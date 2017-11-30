package main

import (
	"Adventure/AdvServer/network"
)

func main() {
	listener := network.NewTCPListenter()
	listener.StartAccept()
}
