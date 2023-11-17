package foliogo


import "os"
import "errors"


func NewDefaultSession() (Session, error) {
	evs := []string{"URL", "TENANT", "USER", "PW"}
	var missing []string

	for _, ev := range(evs) {
		if os.Getenv("OKAPI_" + ev) == "" {
			missing = append(missing, ev)
		}
	}

	if len(missing) > 0 {
		s := "NewDefaultSession: missing environment variables: "
		for i, ev := range(missing) {
			s += "OKAPI_" + ev
			if i < len(missing) - 1 {
				s += ", "
			}
		}
		return Session{}, errors.New(s)
	}

	service := NewService(os.Getenv("OKAPI_URL"))
	return service.Login(os.Getenv("OKAPI_TENANT"), os.Getenv("OKAPI_USER"), os.Getenv("OKAPI_PW"))
}
