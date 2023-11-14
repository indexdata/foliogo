// Package foliogo provides a client library for the FOLIO LMS
package foliogo

import "os"
import "github.com/MikeTaylor/catlogger"


type Service struct {
	url string
	logger *catlogger.Logger
}


func (this Service)String() string {
	return "SERVICE(" + this.url + ")"
}


func NewService(url string) Service {
	logcat := os.Getenv("LOGGING_CATEGORIES")
	if (logcat == "") {
		logcat = os.Getenv("LOGCAT")
	}

	s := Service{
		url: url,
		logger: catlogger.MakeLogger(logcat, "", false),
	}
	s.Log("service", "Made new service on URL", url)
	return s
}


func (this Service)Log(cat string, args ...string) {
	this.logger.Log(cat, args...)
}


func (this Service)Login(tenant string, username string, password string) (Session, error) {
	session := NewSession(this, tenant, username, password)
	err := session.Login()
	if err != nil {
		return Session{}, err
	}
	return session, nil
}
