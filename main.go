package main

import "log"

func main() {
	storage, err := CreateStorage()
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := storage.init(); err != nil {
		log.Fatal(err)
		return
	}
	server := MakeApi(":8080", storage)
	server.Run()
}
