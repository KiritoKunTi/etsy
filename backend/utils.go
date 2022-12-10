package main

type ErrorMessage struct{
	Message string `json:"message"`
	Object interface{} `json:"object"`
}