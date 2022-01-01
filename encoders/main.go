package main

import (
	"bytes"
	"encoding/json"
	"log"
)

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName
	ServiceUpdateURL string
}

type ServiceName string

func main() {
	var r = Registration{
		ServiceName:      ServiceName("adam"),
		ServiceURL:       "localhost:8080",
		RequiredServices: []ServiceName{"adam", "daniel"},
		ServiceUpdateURL: "lala",
	}
	buf := new(bytes.Buffer)
	log.Println(buf)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		log.Println("error while encoding")
	}
	log.Println(buf)
}
