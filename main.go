package main

import "godynamicserver/server"

func main() {
	server := server.NewDServer()
	server.Start()
}
