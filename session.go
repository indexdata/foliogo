package foliogo


import "io"
import "fmt"
import "strings"
import "encoding/json"
import "net/http"
import "net/http/cookiejar"

import "moul.io/http2curl"


// Hash is useful for writing out complex constants for JSON serialization
type Hash map[string]interface{}


type Session struct {
	service Service
	tenant string
	username string
	password string
	client http.Client
}


type RequestParams struct {
	method string
	body string
	json interface{}
}


func (this Session)String() string {
	return "SESSION(" + this.tenant + "/" + this.username + ")"
}


func NewSession(service Service, tenant string, username string, password string) (Session, error) {
	jar, err := cookiejar.New(nil) // XXX or &cookiejar.Options{}
	if err != nil {
		return Session{}, err;
	}
	this := Session{
		service: service,
		tenant: tenant,
		username: username,
		password: password,
		client: http.Client{
			Jar: jar, // Binks
		},
	}

	service.Log("session", "Made new session for service", service.String())
	return this, nil
}


func (this Session)Log(cat string, args ...string) {
	this.service.Log(cat, args...)
}


func (this Session)Fetch(path string, params RequestParams) (Hash, error) {
	var body string
	var err error
	if (params.json != nil) {
		bytes, err := json.Marshal(params.json)
		if err != nil {
			return Hash{}, err
		}
		body = string(bytes)
	} else {
		body = params.body
	}

	method := params.method
	if (method == "") {
		if (body == "") {
			method = "POST"
		} else {
			method = "GET"
		}
	}

	url := this.service.url + "/" + path
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return Hash{}, err
	}
	req.Header.Add("X-Okapi-Tenant", this.tenant)
	this.Log("curl", http2curl.GetCurlCommand(req))

	resp, err := this.client.Do(req)	
	if err != nil {
		return Hash{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, *MakeHttpError(resp.StatusCode, method, url)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s := string(bytes)
	this.Log("body", s)

	return Hash{"error": "not yet implemented"}, nil
}


func (this Session)Login() error {
	this.Log("op", "login(user=" + this.username + ")")
	this.Log("auth", "trying new-style authentication with expiry")
	body := Hash{ "tenant": this.tenant, "username": this.username, "password": this.password }
	json, err := this.Fetch("authn/login-with-expiry", RequestParams{
		json: body,
	})
	if err != nil { return err }
	fmt.Println(json)
	return nil
	/*
	then := Date.parse(json) // XXX should be one element thereof
	now := Date.now()
	ttl := then - now
	fst := process.env.FOLIOJS_SESSION_TIMEOUT
	if (fst == "") {
		timeout = ttl / 2
	} else {
		timeout := fst * 1000
	}
	this.sessionCookie = json.sessionCookie
	// XXX remember to refresh token as needed
	*/
}
