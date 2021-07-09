package main

import (
	"log"
	"net"

	"github.com/kipukun/grtorrent/xmlrpc"
)

type Response struct {
	Raw string `xml:",innerxml"`
}

type Request struct {
}

func main() {
	sock, err := net.Dial("unix", "/home/kipu/rtorrent.socket")
	if err != nil {
		log.Panicln(err)
	}
	c := xmlrpc.NewClient(sock)
	defer c.Close()

	var res Response
	err = c.Call("system.listMethods", "", &res)
	if err != nil {
		log.Panicln(err)
	}
}
