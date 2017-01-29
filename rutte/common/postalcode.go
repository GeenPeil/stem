package common

import "regexp"

//RegexpValidatePostalcode can be used to validate a dutch postalcode (e.g. `1234AB`)
var RegexpValidatePostalcode = regexp.MustCompile(`^[0-9]{4} ?[a-zA-Z]{2}$`)
