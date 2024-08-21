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


func NewService(url string, optLogger ...*catlogger.Logger) Service {
	var logger *catlogger.Logger

	if (len(optLogger) > 0 && optLogger[0] != nil) {
		logger = optLogger[0]
	} else {
		logcat := os.Getenv("LOGGING_CATEGORIES")
		if (logcat == "") {
			logcat = os.Getenv("LOGCAT")
		}
		logger = catlogger.MakeLogger(logcat, "", false)
	}

	s := Service{
		url: url,
		logger: logger,
	}
	s.Log("service", "Made new service on URL", url)
	return s
}


func (this Service)Log(cat string, args ...string) {
	this.logger.Log(cat, args...)
}


func (this Service)Login(tenant string, username string, password string) (Session, error) {
	session, err := NewSession(this, tenant, username, password)
	if err != nil {
		return Session{}, err
	}

	err = session.Login()
	if err != nil {
		return Session{}, err
	}

	return session, nil
}


func (this Service)ResumeSession(tenant string) (Session, error) {
	session, err := NewSession(this, tenant, "", "NOPASS")
	if err != nil {
		return Session{}, err
	}

	return session, nil
}
