package common

import "regexp"

// RegexpValidateEmailAddress is just a very simple regular expression to catch 99% of user input mistakes.
// Email addresses in the system are always validated by actually sending a verification email to it.
var RegexpValidateEmailAddress = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
