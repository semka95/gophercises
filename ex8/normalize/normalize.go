package normalize

import (
	"fmt"
	"unicode"
)

func PhoneNumberNormalize(phone, format string) (string, error) {
	var digits []interface{}
	for _,r := range phone {
		if unicode.IsDigit(r) {
			digits = append(digits, int(r - '0'))
		}
	}

	if len(digits) != 10 {
		return "", fmt.Errorf("phone number \"%s\" should contain 10 digits", phone)
	}

	return fmt.Sprintf(format, digits...), nil
}