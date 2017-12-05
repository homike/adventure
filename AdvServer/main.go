package main

import (
	network "Adventure/AdvServer/network2"
)

func main() {
	listener := network.NewTCPServer()
	listener.Run()
}
