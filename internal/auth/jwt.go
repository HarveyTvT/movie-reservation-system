package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/config"
)

var ErrInvalidToken = errors.New("invalid token")

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func (header JwtHeader) String() string {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(headerBytes)
}

type JwtPayload struct {
	Username string                      `json:"username"`
	Role     movie_reservation.User_Role `json:"role"`
}

func (payload JwtPayload) String() string {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(payloadBytes)
}

type Jwt struct {
	Header  JwtHeader
	Payload JwtPayload
	Sign    string
}

func (jwt Jwt) String() string {
	return jwt.Header.String() + "." + jwt.Payload.String() + "." + jwt.Sign
}

func NewJwtHeader() JwtHeader {
	return JwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}
}

func NewJwtPayload(username string, role movie_reservation.User_Role) JwtPayload {
	return JwtPayload{
		Username: username,
		Role:     role,
	}
}

func NewJwt(payload JwtPayload, secret string) Jwt {
	header := NewJwtHeader()
	sign := Sign(header, payload, secret)

	return Jwt{
		Header:  header,
		Payload: payload,
		Sign:    sign,
	}
}

func Parse(token string) (Jwt, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Jwt{}, ErrInvalidToken
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	sign := parts[2]

	headerBytes, err := base64.StdEncoding.DecodeString(headerEncoded)
	if err != nil {
		return Jwt{}, ErrInvalidToken
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return Jwt{}, ErrInvalidToken
	}

	header := JwtHeader{}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return Jwt{}, ErrInvalidToken
	}

	payload := JwtPayload{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return Jwt{}, ErrInvalidToken
	}

	if Sign(header, payload, config.Get().Secret) != sign {
		return Jwt{}, ErrInvalidToken
	}

	return Jwt{
		Header:  header,
		Payload: payload,
		Sign:    sign,
	}, nil

}

func Sign(header JwtHeader, payload JwtPayload, secret string) string {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return ""
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	headerEncoded := base64.StdEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.StdEncoding.EncodeToString(payloadBytes)

	signer := hmac.New(sha256.New, []byte(secret))
	signer.Write([]byte(headerEncoded + "." + payloadEncoded))
	sign := base64.StdEncoding.EncodeToString(signer.Sum(nil))

	return sign
}
