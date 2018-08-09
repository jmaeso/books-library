package api

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
)

func getStringFromSession(r *http.Request, key string) string {
	var strVal string
	if val := sessions.GetSession(r).Get(key); val != nil {
		strVal = val.(string)
	}

	return strVal
}
