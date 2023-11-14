package main

import "fmt"
import "github.com/indexdata/foliogo"

func main() {
	service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org")
	fmt.Printf("got service %+v\n", service)
	session := service.Login("diku", "user-basic-view", "user-basic-view")
	fmt.Printf("got session %+v\n", session)
	/*
	body := session.folioFetch("/users?limit=20")
	fmt.Println(body.users)
	session.close()
	*/
}
