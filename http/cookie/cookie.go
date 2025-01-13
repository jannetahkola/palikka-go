package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	palikka "palikka-go"
	"time"
)

// defaults holds static values for http.Cookie
// that cannot be provided with Options.
var defaults = struct {
	HttpOnly bool
	Path     string
	SameSite http.SameSite
}{
	HttpOnly: true,
	Path:     "/",
	SameSite: http.SameSiteStrictMode,
}

type Options struct {
	CookieName string
	Secret     []byte
	Secure     bool
}

type Service interface {
	Parse(r *http.Request) (*http.Cookie, error)
	New(session *palikka.Session) *http.Cookie
	NewExpired() *http.Cookie
	Verify(c *http.Cookie) (string, error)
	CookieName() string
}

func NewService(options *Options) Service {
	return &serviceImpl{options: options}
}

type serviceImpl struct {
	options *Options
}

func (s *serviceImpl) CookieName() string {
	return s.options.CookieName
}

func (s *serviceImpl) Parse(r *http.Request) (*http.Cookie, error) {
	sessionCookie, _ := r.Cookie(s.options.CookieName)

	if sessionCookie == nil {
		return nil, errors.New("no cookie in request")
	}

	return sessionCookie, nil
}

func (s *serviceImpl) New(session *palikka.Session) *http.Cookie {

	// todo maybe document the whole thing in readme.md or something
	// NOTE: required if using 'localhost' when developing:
	// - HttpOnly: true
	// - Secure: false
	// - Path: "/"

	cookie := &http.Cookie{
		Name:     s.options.CookieName,
		Value:    session.ID,
		MaxAge:   int(session.MaxAge.Seconds()),
		Expires:  session.Expires,
		Secure:   false,
		HttpOnly: defaults.HttpOnly,
		Path:     defaults.Path,
		SameSite: defaults.SameSite,
	}

	hash := hmac.New(sha256.New, s.options.Secret)
	hash.Write([]byte(cookie.Name))
	hash.Write([]byte(cookie.Value))
	cookie.Value = hex.EncodeToString(hash.Sum(nil)) + cookie.Value // todo hide the SessionID?

	return cookie
}

func (s *serviceImpl) Verify(cookie *http.Cookie) (string, error) {
	encodedSignatureSize := sha256.Size * 2

	if len(cookie.Value) < encodedSignatureSize {
		return "", errors.New(fmt.Sprintf(
			"invalid cookie value, size=%d, expected=%d", len(cookie.Value), encodedSignatureSize))
	}

	cookieSignature, _ := hex.DecodeString(cookie.Value[:encodedSignatureSize])
	cookieValue := cookie.Value[encodedSignatureSize:]

	hash := hmac.New(sha256.New, s.options.Secret)
	hash.Write([]byte(cookie.Name))
	hash.Write([]byte(cookieValue))
	expectedSignature := hash.Sum(nil)

	if !hmac.Equal(cookieSignature, expectedSignature) {
		return "", errors.New("invalid cookie signature")
	}

	return cookieValue, nil
}

func (s *serviceImpl) NewExpired() *http.Cookie {
	return &http.Cookie{
		Name:     s.options.CookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Secure:   s.options.Secure,
		HttpOnly: defaults.HttpOnly,
		Path:     defaults.Path,
		SameSite: defaults.SameSite,
	}
}
