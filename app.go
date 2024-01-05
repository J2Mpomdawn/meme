package main

import (
	"embed"

	"meme/control"
)

//go:embed views/templates views/scripts views/styles views/images/icon
var resources embed.FS

func main() {
	control.Run(resources)
}
