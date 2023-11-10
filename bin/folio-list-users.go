package main

import "os"
import "fmt"
import "github.com/indexdata/foliogo"

func main() {
	service := Folio.service('https://folio-snapshot-okapi.dev.folio.org')
	session := service.login('diku', 'user-basic-view', 'user-basic-view')
        body := session.folioFetch('/users?limit=20')
	fmt.Println(body.users)
        session.close()
}


