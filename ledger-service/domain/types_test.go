package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	t.Run("creates currency from string", func(t *testing.T) {
		c := Currency("USD")
		assert.Equal(t, Currency("USD"), c)
	})

	t.Run("empty currency is valid", func(t *testing.T) {
		c := Currency("")
		assert.Equal(t, Currency(""), c)
	})
}

func TestAccountStatus(t *testing.T) {
	t.Run("creates status from string", func(t *testing.T) {
		s := AccountStatus("active")
		assert.Equal(t, AccountStatus("active"), s)
	})

	t.Run("empty status is valid", func(t *testing.T) {
		s := AccountStatus("")
		assert.Equal(t, AccountStatus(""), s)
	})

	t.Run("common status values", func(t *testing.T) {
		assert.Equal(t, AccountStatus("active"), AccountStatus("active"))
		assert.Equal(t, AccountStatus("inactive"), AccountStatus("inactive"))
		assert.Equal(t, AccountStatus("frozen"), AccountStatus("frozen"))
	})
}
