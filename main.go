package main

import (
	"log"
	"toolset/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Executr err: %v",err)
	}
}
