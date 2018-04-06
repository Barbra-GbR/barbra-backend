package helpers

import "regexp"

var (
	RegexEmail    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	RegexName     = regexp.MustCompile("^[\\p{L}.-]+$")
	RegexURL      = regexp.MustCompile("^https?:\\/\\/[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b[-a-zA-Z0-9@:%_\\+.~#?&\\/\\/=]*$")
	RegexNickname = regexp.MustCompile("^[\\p{L}-\\._0-9]+$")
)
