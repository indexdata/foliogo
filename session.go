package foliogo


type Session struct {
	service Service
	tenant string
	username string
	password string
}


func NewSession(service Service, tenant string, username string, password string) Session {
	this := Session{
		service: service,
		tenant: tenant,
		username: username,
		password: password,
	}
	service.Log("session", "Made new session for service", service.String());
	return this;
}


func (this Session)login() {
	// XXX to do
}
