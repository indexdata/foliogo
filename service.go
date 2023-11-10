package foliogo

import "os"
import "github.com/MikeTaylor/catlogger"


type Service struct {
	url string
	logger *catlogger.Logger
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
	s.Log("service", "Made new service on URL", url);
	return s
}

func (s Service)Log(cat string, args ...string) {
	s.logger.Log(cat, args...)
}

/*
  async login(tenant, username, password) {
    const session = new FolioSession(this, tenant, username, password);
    await session.login();
    return session;
  }

  resumeSession(tenant, token) {
    const session = new FolioSession(this, tenant);
    session.token = token; // It might be more polite to use an API here
    return session;
  }
}
*/
