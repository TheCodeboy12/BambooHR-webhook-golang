package helpers

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
)

func ValidateHmac(payload string, signature string, secret string) bool {

	var bSignature []byte
	if IsHex(signature) {
		bSignature, _ = hex.DecodeString(signature)
	} else if IsBase64(signature) {
		bSignature, _ = base64.StdEncoding.DecodeString(signature)
	} else {
		slog.Error("Invalid signature format")
		return false
	}

	mac := hmac.New(crypto.SHA256.New, []byte(secret))
	_, err := mac.Write([]byte(payload))
	if err != nil {
		slog.Error("Error writing to mac", "error", err.Error())
		return false
	}
	expectedMAC := mac.Sum(nil)
	ret := hmac.Equal(expectedMAC, bSignature)
	return ret

}
