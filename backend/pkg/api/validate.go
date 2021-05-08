package api

import "regexp"

var ptnCommitID *regexp.Regexp

func initValidate() {
	ptnCommitID = regexp.MustCompile("^[0-9a-f]{40}$")
}

func isValidCommitID(s string) bool {
	return ptnCommitID.MatchString(s)
}
