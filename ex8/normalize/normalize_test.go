package normalize

import "testing"

func TestPhoneNumberNormalize(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		template string
		want     string
		errMsg   string
	}{
		{
			"all digits",
			"1234567890",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567890",
			"",
		},
		{
			"digits with spaces",
			"123 456 7891",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567891",
			"",
		},
		{
			"digits with spaces and parentheses",
			"(123) 456 7892",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567892",
			"",
		},
		{
			"digits with spaces, parentheses and dash",
			"(123) 456-7893",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567893",
			"",
		},
		{
			"digits dashes",
			"123-456-7894",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567894",
			"",
		},
		{
			"digits, parentheses and dashes",
			"(123)456-7892",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"(123)4567892",
			"",
		},
		{
			"bad phone",
			"(13)456-7892",
			"(%d%d%d)%d%d%d%d%d%d%d",
			"",
			"phone number \"(13)456-7892\" should contain 10 digits",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := PhoneNumberNormalize(test.phone, test.template)
			if test.errMsg != "" {
				if err.Error() != test.errMsg && got == "" {
					t.Errorf("should be \"%s\" error, but get \"%s\"", test.errMsg, err.Error())
				}
			}

			if got != test.want {
				t.Errorf("got %s, wanted %s", got, test.want)
			}

		})
	}
}
