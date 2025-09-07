package encoder

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// EncodeBase64 returns base64-encoded string of input bytes
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 decodes base64 string into plain text
func DecodeBase64(input string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// EncodeHex returns hex-encoded string of input bytes
func EncodeHex(input string) string {
	return hex.EncodeToString([]byte(input))
}

// DecodeHex decodes hex string into plain text
func DecodeHex(input string) (string, error) {
	b, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// EncodeURL performs URL percent-encoding
func EncodeURL(input string) string {
	return url.QueryEscape(input)
}

// DecodeURL performs URL percent-decoding
func DecodeURL(input string) (string, error) {
	return url.QueryUnescape(input)
}
