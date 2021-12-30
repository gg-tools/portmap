package internal

import "github.com/gg-tools/portmap/internal/link"

func Serve() {
	server := link.NewServer("123")
	link.Listen(":9988", server.OnConnect)
	link.Listen(":10022", server.OnUserConnect)
}
