package input_test

import (
	"github.com/stretchr/testify/assert"
	"palikka-go/http/input"
	"testing"
)

func TestInput_Sanitize(t *testing.T) {
	t.Run("removes spaces", func(t *testing.T) {
		v := " password \n\t"
		result := input.Sanitize(&v)
		assert.Equal(t, "password", *result)
	})

	t.Run("returns nil if value is nil", func(t *testing.T) {
		result := input.Sanitize(nil)
		assert.Nil(t, result)
	})
}

func TestInput_ValidateUsername(t *testing.T) {
	t.Run("returns error if value is nil", func(t *testing.T) {
		err := input.ValidateUsername(nil)
		assert.ErrorContains(t, err, "username is empty")
	})

	t.Run("returns error if value is empty", func(t *testing.T) {
		v := ""
		err := input.ValidateUsername(&v)
		assert.ErrorContains(t, err, "username is empty")
	})

	t.Run("returns error if value is too short", func(t *testing.T) {
		v := "user"
		err := input.ValidateUsername(&v)
		assert.ErrorContains(t, err, "username is invalid, must be between")
	})

	t.Run("returns error if value is too long", func(t *testing.T) {
		v := "usernameusernameusernameu"
		err := input.ValidateUsername(&v)
		assert.ErrorContains(t, err, "username is invalid, must be between")
	})

	t.Run("returns error if value contains invalid chars", func(t *testing.T) {
		v := []string{
			"username@",
			"user name",
			"username,",
		}

		for _, username := range v {
			err := input.ValidateUsername(&username)
			assert.ErrorContains(t, err, "username is invalid, must match")
		}
	})

	t.Run("returns ok", func(t *testing.T) {
		v := "123username-_"
		err := input.ValidateUsername(&v)
		assert.NoError(t, err)
	})
}

func TestInput_ValidatePassword(t *testing.T) {
	t.Run("returns error if value is nil", func(t *testing.T) {
		err := input.ValidatePassword(nil)
		assert.ErrorContains(t, err, "password is empty")
	})

	t.Run("returns error if value is empty", func(t *testing.T) {
		v := ""
		err := input.ValidatePassword(&v)
		assert.ErrorContains(t, err, "password is empty")
	})

	t.Run("returns error if value is too short", func(t *testing.T) {
		v := "pass"
		err := input.ValidatePassword(&v)
		assert.ErrorContains(t, err, "password is invalid, must be between")
	})

	t.Run("returns error if value is too long", func(t *testing.T) {
		v := ""
		for i := 0; i < 72+1; i++ {
			v += "p"
		}
		err := input.ValidatePassword(&v)
		assert.ErrorContains(t, err, "password is invalid, must be between")
	})

	t.Run("returns ok", func(t *testing.T) {
		v := "123passw ord-_@.£$<\t\b "
		err := input.ValidatePassword(&v)
		assert.NoError(t, err)
	})
}
