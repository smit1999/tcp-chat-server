package main

type room struct {
	name    string
	members map[string]*client
}
