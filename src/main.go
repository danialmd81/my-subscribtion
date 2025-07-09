package main

import (
	"github.com/danialmd81/my-subscribtion/subs"
	"github.com/danialmd81/my-subscribtion/telegram"
)

func main() {
	telegram.Run()
	subs.Run()
}
