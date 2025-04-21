package foliogo


import "os"
import "io"
import "strings"
import "time"
import "strconv"
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
	accessToken string // Used only for "curl"-category logging
}


type RequestParams struct {
	Method string
	Body string
	ContentType string
	Json interface{}
	Token string
}


func (this Session)String() string {
	return "SESSION(" + this.tenant + "/" + this.username + ":" + this.refreshAfter.String() + ")"
}


func NewSession(service Service, tenant string, username string, password string) (Session, error) {
	jar, err := cookiejar.New(nil) // XXX or &cookiejar.Options{}
	if err != nil {
		return Session{}, err
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


func (this Session)GetTenant() string {
	return this.tenant
}

func (this Session)Log(cat string, args ...string) {
	this.service.Log(cat, args...)
}


// Private to Fetch
func extractAccessToken(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "folioAccessToken" {
			return cookie.Value
			break
		}
	}
	return ""
}


func (this *Session)Fetch(path string, params RequestParams) ([]byte, error) {
	// Check whether a renewed login is required
	if (!this.refreshAfter.IsZero() &&
		time.Now().Compare(this.refreshAfter) > 0) {
		this.Log("auth", "refresh authentication")
		this.refreshAfter = time.Time{} // The zero value, prevents recursion loop
		err := this.Login()
		if err != nil {
			return []byte{}, err
		}
	}

	var body string
	var err error
	if (params.Json != nil) {
		bytes, err2 := json.Marshal(params.Json)
		if err2 != nil {
			return []byte{}, err
		}
		body = string(bytes)
	} else {
		body = params.Body
	}

	method := params.Method
	if (method == "") {
		if (body == "") {
			method = "GET"
		} else {
			method = "POST"
		}
	}

	url := this.service.url + "/" + path
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("X-Okapi-Tenant", this.tenant)
	if params.ContentType != "" {
		req.Header.Add("Content-type", params.ContentType)
	} else if params.Json != nil {
		req.Header.Add("Content-type", "application/json")
	}
	if params.Token != "" {
		req.Header.Add("X-Okapi-Token", params.Token)
	}

	curlCommand, _ := http2curl.GetCurlCommand(req)
	curlString := curlCommand.String()
	if this.accessToken != "" && params.Token == "" {
		curlString = strings.Replace(curlString, " ", " -H 'X-Okapi-Token: " + this.accessToken + "' ", 1)
	}
	this.Log("curl", curlString)
	resp, err := this.client.Do(req)
	if err != nil {
		// I think this is for a low-level error such as DNS resolution failure
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Special case: on login, remember access token for subsequent "curl" logging
	if path == "authn/login-with-expiry" {
		this.accessToken = extractAccessToken(resp)
	}

	contentType := resp.Header.Get("Content-Type")
	this.Log("status", resp.Status, "(" + contentType + ")")

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	this.Log("response", string(bytes))
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, *MakeHTTPError(resp.StatusCode, method, url, string(bytes))
	}

	return bytes, nil
}


func (this *Session)Fetch0(path string) ([]byte, error) {
	return this.Fetch(path, RequestParams{})
}


func (this *Session)Login() error {
	this.Log("op", "login(user=" + this.username + ")")
	this.Log("auth", "trying new-style authentication with expiry")
	body := Hash{ "tenant": this.tenant, "username": this.username, "password": this.password }
	bytes, err := this.Fetch("authn/login-with-expiry", RequestParams{
		Json: body,
	})
	if err != nil {
		return err
	}

	timeout := os.Getenv("FOLIOGO_SESSION_TIMEOUT")
	if timeout != "" {
		// No need to consult the HTTP response at all!
		secs, err2 := strconv.Atoi(timeout)
		if err2 != nil {
			return err2
		}
		this.refreshAfter = time.Now().Add(time.Duration(secs) * time.Second)
		return nil
	}

	var data Hash
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	expiryString := data["accessTokenExpiration"].(string)
	// We don't really care about refreshTokenExpiration
	expiryTime, err :=  time.Parse(time.RFC3339, expiryString)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	diff := expiryTime.Sub(now)
	diff90 := 9 * diff / 10
	this.refreshAfter = now.Add(diff90)

	return nil
}
