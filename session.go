package foliogo


import "io"
import "fmt"
import "strings"
import "time"
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
	refreshAfter time.Time
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
			method = "GET"
		} else {
			method = "POST"
		}
	}

	url := this.service.url + "/" + path
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return Hash{}, err
	}
	req.Header.Add("X-Okapi-Tenant", this.tenant)
	curlCommand, _ := http2curl.GetCurlCommand(req)
	// I don't know why I need this Sprintf, but curlCommand is not a string and cannot be simply converted
	this.Log("curl", fmt.Sprintf("%s", curlCommand))

	resp, err := this.client.Do(req)	
	if err != nil {
		return Hash{}, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	this.Log("status", resp.Status, "(" + contentType + ")")

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, *MakeHttpError(resp.StatusCode, method, url)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	this.Log("response", string(bytes))

	// Every valid FOLIO WSAPI is JSON
	var data Hash
	err = json.Unmarshal(bytes, &data)
	return data, nil
}


func (this Session)Login() error {
	this.Log("op", "login(user=" + this.username + ")")
	this.Log("auth", "trying new-style authentication with expiry")
	body := Hash{ "tenant": this.tenant, "username": this.username, "password": this.password }
	data, err := this.Fetch("authn/login-with-expiry", RequestParams{
		json: body,
	})
	if err != nil {
		return err
	}

	expiryString := data["accessTokenExpiration"].(string)
	// We don't really care about refreshTokenExpiration
	expiryTime, err :=  time.Parse(time.RFC3339, expiryString)
	now := time.Now().UTC()
	diff := expiryTime.Sub(now)
	diff90 := 9 * diff / 10
	this.refreshAfter = now.Add(diff90);

	return nil
}
