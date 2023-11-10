package main

import "fmt"
import "github.com/indexdata/foliogo"

func main() {
	service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org")
	fmt.Printf("got service %+v", service)
	/*
	session := service.login("diku", "user-basic-view", "user-basic-view")
	body := session.folioFetch("/users?limit=20")
	fmt.Println(body.users)
	session.close()
	*/
}
