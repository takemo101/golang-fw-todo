package main

import (
	"app/chiapp"
	"app/echoapp"
	"app/fiberapp"
	"app/ginapp"
	"app/shared"
	"flag"
	"log"
)

var ServerTypes = map[string]func() shared.Server{
	"gin":   ginapp.NewGinServer,
	"fiber": fiberapp.NewFiberServer,
	"echo":  echoapp.NewEchoServer,
	"chi":   chiapp.NewChiServer,
}

// ServerTypeに対応したサーバーを生成する
func newServer(
	t string,
) (shared.Server, bool) {

	if fn, ok := ServerTypes[t]; ok {
		return fn(), true
	}

	return nil, false
}

func main() {
	// コマンド引数からサーバータイプを取得
	// go run main.go -server fiber
	t := flag.String("server", "gin", "server type")
	host := flag.String("host", "localhost:8080", "host")
	flag.Parse()

	if server, ok := newServer(*t); ok {
		server.Run(*host)
	} else {
		log.Fatal("invalid server type")
	}
}
