package main

import "os"
import "fmt"
import "time"
import "strconv"
import "encoding/json"
import "github.com/indexdata/foliogo"
import "github.com/MikeTaylor/catlogger"

type user struct {
	Username string `json:"username"`
	Active bool `json:"active"`
}
type response struct {
	Users []user `json:"users"`
	TotalRecords int `json:"totalRecords"y`
}

func main() {
	var logger *catlogger.Logger
	if (false) {
		logger = catlogger.MakeLogger(os.Getenv("LOGCAT"), "", true)
	}
	var session foliogo.Session
	if (true) {
		var err error
		service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org", logger)
		session, err = service.Login("diku", "user-basic-view", "user-basic-view")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: login failed: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("logger = %x\n", logger);
		// Run with (for example): OKAPI_URL=https://folio-snapshot-okapi.dev.folio.org OKAPI_TENANT=diku OKAPI_USER=user-basic-view OKAPI_PW=user-basic-view
		var err error
		session, err = foliogo.NewDefaultSession(logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: default session failed: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 {
		nsecs, _ := strconv.Atoi(os.Args[1])
		time.Sleep(time.Duration(nsecs) * time.Second)
	}

	bytes, err := session.Fetch0("users?limit=20")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: fetch users failed: %s\n", os.Args[0], err)
		os.Exit(2)
	}

	var r response
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: decode JSON failed: %s\n", os.Args[0], err)
		os.Exit(2)
	}

	users := r.Users
	for _, user := range users {
		var marker string
		if user.Active {
			marker = "*"
		} else {
			marker = " "
		}
		fmt.Println(marker, user.Username)
	}
}
