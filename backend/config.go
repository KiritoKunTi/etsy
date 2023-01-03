package main

import "flag"

type Config struct {
	Address string
	Static  string
	Private string
}

var config Config

func init() {
	config = Config{
		Address: *flag.String("address", "127.0.0.1:8080", "ip address and port number of our website"),
		Static:  *flag.String("static", "static", "static files location of our website"),
		Private: *flag.String("private", "private", "private files location of our website"),
	}
}
