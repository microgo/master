package utils

import (
	"github.com/srinathh/hashtag"
	"strings"
)

func ExtractHashtags(text string) []string {
	str := strings.Replace(text, "#", " #", -1)
	return hashtag.ExtractHashtags(str)
}
