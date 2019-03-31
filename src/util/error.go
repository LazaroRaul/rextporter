package util

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// ErrorFromThisScope can be used to get notified about a general error(what do you want from this function)
// and a specific error(`rootCause`), the statement who causes the error
// this method will log the rootCause and return the general error(`generalScopeErr`)
func ErrorFromThisScope(rootCause, generalScopeErr string) error {
	log.WithError(errors.New(rootCause)).Errorln("root cause error")
	return errors.New(generalScopeErr)
}
