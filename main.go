package main

import (
	"github.com/axetroy/go-server/src"
)

func main() {
	go src.ServerUserClient()
	src.ServerAdminClient()
}
