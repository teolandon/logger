package logger

import (
	"crypto/md5"
)

var containsEmpty = false

type set map[[16]byte]string

func (s set) add(str string) {
	if str == "" {
		containsEmpty = true
		return
	}

	hash := md5.Sum([]byte(str))
	s[hash] = str
}

func (s set) remove(str string) {
	if str == "" {
		containsEmpty = false
		return
	}

	hash := md5.Sum([]byte(str))
	s[hash] = ""
}

func (s set) contains(str string) bool {
	if str == "" {
		return containsEmpty
	}

	hash := md5.Sum([]byte(str))
	return s[hash] == str
}
