package foliogo

import "fmt"
import "net/http"

type httpError struct {
	status int
	method string
	url string
}

func MakeHTTPError(status int, 	method string, url string) *httpError {
	return &httpError{status: status, method: method, url: url}
}

func (this httpError) String() string {
	return fmt.Sprintf("HTTP error %d (%s): %s %s", this.status, http.StatusText(this.status), this.method, this.url)
}

func (this httpError) Error() string {
	return this.String()
}
