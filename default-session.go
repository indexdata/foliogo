package foliogo


import "os"
import "errors"
import "github.com/MikeTaylor/catlogger"


func NewDefaultSession(optLogger ...*catlogger.Logger) (Session, error) {
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

	var service Service
	if (len(optLogger) > 0 && optLogger[0] != nil) {
		service = NewService(os.Getenv("OKAPI_URL"), optLogger[0])
	} else {
		service = NewService(os.Getenv("OKAPI_URL"))
	}

	return service.Login(os.Getenv("OKAPI_TENANT"), os.Getenv("OKAPI_USER"), os.Getenv("OKAPI_PW"))
}
