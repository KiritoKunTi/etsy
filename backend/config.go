package main

import "flag"

type Config struct{
	Address string
	Static string
	Private string
}

var config Config

func init(){
	flag.Parse()
	config = Config{
		Address: *flag.String("address", "127.0.0.1:8080", "ip address and port number of our website"),
		Static: *flag.String("static", "127.0.0.1", "static files location of our website"),
		Private: *flag.String("address", "127.0.0.1", "private files location of our website"),
	}
}