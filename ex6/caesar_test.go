package ex6

import "testing"

func TestCaesarCipher(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		inputShift  int
		want        string
	}{
		{
			"8 words",
			"Always-Look-on-the-Bright-Side-of-Life",
			5,
			"Fqbfdx-Qttp-ts-ymj-Gwnlmy-Xnij-tk-Qnkj",
		},
		{
			"2 words",
			"Hello_World!",
			4,
			"Lipps_Asvph!",
		},
		{
			"26 shift, nothing changed",
			"Ciphering.",
			26,
			"Ciphering.",
		},
		{
			"url",
			"www.abc.xy",
			87,
			"fff.jkl.gh",
		},
		{
			"random",
			"159357lcfd",
			98,
			"159357fwzx",
		},
		{
			"0 shift",
			"D3q4",
			0,
			"D3q4",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := caesarCipher(test.inputString, test.inputShift)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}
