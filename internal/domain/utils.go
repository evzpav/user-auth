package domain

import "regexp"

var rxEmail = regexp.MustCompile(".+@.+\\..+")

func validateEmail(email string) bool {
	return rxEmail.Match([]byte(email))
}
