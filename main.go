package main

import (
	"log"

	"gitlab.com/middlefront/workspace/cmd"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.Execute()
}
