package utils

import "regexp"

var (
	HttpPattern = regexp.MustCompile(`(?i)^(GET|POST|PUT|DELETE|OPTIONS|HEAD|PATCH)`)
	SshPattern  = regexp.MustCompile(`^(SSH-)`)
)
