package main

import "os"
import "fmt"
import "github.com/indexdata/foliogo"

func main() {
	service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org")
	fmt.Printf("got service %+v\n", service)
	session, err := service.Login("diku", "user-basic-view", "user-basic-view")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: login failed: %s\n", os.Args[0], err)
		os.Exit(1)
	}
	fmt.Printf("got session %+v\n", session)
	/*
	body := session.folioFetch("/users?limit=20")
	fmt.Println(body.users)
	session.close()
	*/
}
