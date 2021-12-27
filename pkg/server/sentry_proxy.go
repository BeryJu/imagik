package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
)

type SentryRequest struct {
	DSN string `json:"dsn"`
}

func (s *Server) APISentryProxy(rw http.ResponseWriter, r *http.Request) {
	fullBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Debug("failed to read body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	lines := strings.Split(string(fullBody), "\n")
	if len(lines) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	sd := SentryRequest{}
	err = json.Unmarshal([]byte(lines[0]), &sd)
	if err != nil {
		s.logger.WithError(err).Warning("failed to parse sentry request")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	ourDSN := sentry.CurrentHub().Client().Options().Dsn
	if sd.DSN != ourDSN {
		s.logger.WithField("have", sd.DSN).WithField("expected", ourDSN).Debug("invalid DSN")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := http.DefaultClient.Post("https://sentry.beryju.org/api/8/envelope/", "application/octet-stream", strings.NewReader(string(fullBody)))
	if err != nil {
		s.logger.WithError(err).Warning("failed to proxy sentry")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(res.StatusCode)
}
