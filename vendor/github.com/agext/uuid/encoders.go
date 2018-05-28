package uuid

import "encoding/base64"

// Encoder implementations provide a method of encoding a UUID into a byte slice.
type Encoder interface {
	Encode([]byte) []byte
}

// EncoderToString implementations provide a method of encoding a UUID into a string.
type EncoderToString interface {
	EncodeToString([]byte) string
}

var (
	// Base64URLEncoder uses Base64 URL Encoding
	Base64URLEncoder = Base64Encoder{base64.RawURLEncoding}
	// Base64StdEncoder uses Base64 Std Encoding
	Base64StdEncoder = Base64Encoder{base64.RawStdEncoding}
)

// Base64Encoder is a wrapper around any encoding/base64.Encoding to satisfy Encoder and EncoderToString.
type Base64Encoder struct {
	Enc *base64.Encoding
}

// Encode encodes the source to a byte slice using the encoding/base64.Encoding set on the receiver.
func (e Base64Encoder) Encode(src []byte) (out []byte) {
	out = make([]byte, e.Enc.EncodedLen(len(src)))
	e.Enc.Encode(out, src)
	return
}

// EncodeToString encodes the source to a string using the encoding/base64.Encoding set on the receiver.
func (e Base64Encoder) EncodeToString(src []byte) (out string) {
	return string(e.Encode(src))
}
