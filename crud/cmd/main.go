package main

import "ChatService/crud/internal/config"

func main() {
	cnf := config.MustLoad()
	_ = cnf
	//TODO: init logger
	//TODO: init app
	//TODO: graceful shootDawn
}
