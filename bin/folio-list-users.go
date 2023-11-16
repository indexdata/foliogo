package main

import "os"
import "fmt"
import "time"
import "strconv"
import "github.com/indexdata/foliogo"

func main() {
	service := foliogo.NewService("https://folio-snapshot-okapi.dev.folio.org")
	session, err := service.Login("diku", "user-basic-view", "user-basic-view")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: login failed: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		nsecs, _ := strconv.Atoi(os.Args[1])
		time.Sleep(time.Duration(nsecs) * time.Second)
	}

	body, err := session.Fetch("users?limit=20", foliogo.RequestParams{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: fetch users failed: %s\n", os.Args[0], err)
		os.Exit(2)
	}

	// Unfortunately, the rendering part is clumsy in Go
	users := body["users"].([]interface{})
	for _, e := range users {
		user := e.(map[string]interface{})
		var marker string
		if user["active"] == true {
			marker = "*"
		} else {
			marker = " "
		}
		fmt.Println(marker, user["username"])
	}
}
