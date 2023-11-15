package foliogo

import "fmt"
import "net/http"

type httpError struct {
	status int
	method string
	url string
	addInfo string
}

func MakeHTTPError(status int, 	method string, url string, addInfo string) *httpError {
	return &httpError{status: status, method: method, url: url, addInfo: addInfo}
}

func (this httpError) String() string {
	s := fmt.Sprintf("HTTP error %d (%s): %s %s", this.status, http.StatusText(this.status), this.method, this.url)
	if this.addInfo != "" {
		s += " -- " + this.addInfo
	}
	return s
}

func (this httpError) Error() string {
	return this.String()
}
