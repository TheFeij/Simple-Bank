package token

import (
	"Simple-Bank/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomUsername()
	duration := 1 * time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomUsername()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, err, ErrExpiredToken)
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	username := util.RandomUsername()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	require.NotEmpty(t, jwtToken)

	invalidToken, err := jwtToken.SignedString([]byte(util.RandomString(32, util.ALL)))
	require.NoError(t, err)
	require.NotEmpty(t, invalidToken)

	maker, err := NewPasetoMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.Equal(t, ErrInvalidToken, err)
	require.Nil(t, payload)
}
