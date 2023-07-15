package censor

import (
	"fmt"
	"regexp"
)

func Censor(censorWords []string) func(body string) string {
	return func(body string) string {
		for _, word := range censorWords {
			re := regexp.MustCompile(`"` + word + `":\s*"[^"]*"`)
			fmt.Print(re)
			body = re.ReplaceAllString(body, `"`+word+`": "********"`)
			fmt.Print(body)
		}
		return body
	}
}
