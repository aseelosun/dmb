package utils

import "regexp"

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[A-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	fcbDomain := "@1cb.kz"
	if emailRegex.MatchString(e) && (e[len(e)-7:] == fcbDomain) {
		return true
	}
	return false
}
func IsPhoneNumberValid(e string) bool {
	if e[0] == '7' && len(e) == 11 {
		return true
	} else {
		return false
	}
}
