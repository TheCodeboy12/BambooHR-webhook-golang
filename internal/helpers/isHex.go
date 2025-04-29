package helpers

import (
	"encoding/hex"
	"unicode"
)

func IsHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil && allHexCharacters(s) //Ensure only hex characters
}

func allHexCharacters(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return false
		}
	}

	return true
}
