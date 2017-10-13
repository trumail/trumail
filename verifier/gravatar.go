package verifier

import "net/http"

// HasGravatar performs an http HEAD request to check if the email is
// associated with a gravatar account
func HasGravatar(a *Address) bool {
	resp, err := http.Head("https://en.gravatar.com/" + a.MD5() + ".json")
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()
	return true
}
