package cookie_test

import (
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	palikka "palikka-go"
	palikkacookie "palikka-go/http/cookie"
	"testing"
	"time"
)

func TestCookieService(t *testing.T) {
	var (
		cookieName string
		cookie     palikkacookie.Service
		session    *palikka.Session
	)

	cookieName = "mock-cookie-name"
	cookie = palikkacookie.NewService(
		&palikkacookie.Options{
			CookieName: cookieName,
			Secret:     []byte("test-secret"),
		},
	)
	session = &palikka.Session{
		ID:      "mock-session-id",
		UserID:  1,
		MaxAge:  time.Minute,
		Expires: time.Now().Add(time.Minute),
	}

	t.Run("creates new cookie", func(t *testing.T) {
		c := cookie.New(session)
		assert.NotNil(t, c)
		assert.Equal(t, cookieName, c.Name)
		assert.Equal(t, int(time.Minute.Seconds()), c.MaxAge)
		assert.True(t, c.Expires.Before(time.Now().Add(time.Minute*2)))
		assert.True(t, c.Expires.After(time.Now().Add(time.Minute*-1)))
		assert.True(t, c.Secure)
		assert.True(t, c.HttpOnly)

		// check value is signed
		assert.True(t, len(c.Value) >= sha256.Size*2)
		assert.NotEqual(t, "mock-session-id", c.Value)
		assert.Contains(t, c.Value, "mock-session-id")
	})

	t.Run("creates new expired cookie", func(t *testing.T) {
		c := cookie.NewExpired()
		assert.NotNil(t, c)
		assert.Equal(t, cookieName, c.Name)
		assert.Empty(t, c.Value)
		assert.Equal(t, -1, c.MaxAge)
		assert.True(t, c.Expires.Before(time.Now()))
	})

	t.Run("verifies cookie", func(t *testing.T) {
		c := cookie.New(session)
		v, err := cookie.Verify(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
		assert.Equal(t, "mock-session-id", v)
	})

	t.Run("verifies expired cookie", func(t *testing.T) {
		c := cookie.NewExpired()
		v, err := cookie.Verify(c)
		assert.Error(t, err)
		assert.Empty(t, v)
	})

	t.Run("verifies cookie with invalid signature", func(t *testing.T) {
		cookie2 := palikkacookie.NewService(
			&palikkacookie.Options{
				CookieName: cookieName,
				Secret:     []byte("invalid-secret"),
			},
		)

		c := cookie2.New(session)

		_, err := cookie.Verify(c)
		assert.Error(t, err)
		assert.Equal(t, "invalid cookie signature", err.Error())
	})

	t.Run("parses cookie", func(t *testing.T) {
		c := cookie.New(session)

		req := httptest.NewRequest("GET", "/", nil)

		cc, err := cookie.Parse(req)
		assert.Error(t, err)
		assert.Equal(t, "no cookie in request", err.Error())

		req.AddCookie(c)

		cc, err = cookie.Parse(req)
		assert.NoError(t, err)
		assert.NotNil(t, cc)
	})
}
