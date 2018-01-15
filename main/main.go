package main

import "github.com/teolandon/logger"

func main() {
	logger.Init("logger")
	logger.Println("testing...")
}
