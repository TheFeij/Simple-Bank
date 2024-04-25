package token

import (
	"Simple-Bank/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomUsername()
	duration := 1 * time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	returnedPayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, returnedPayload)

	require.Equal(t, username, returnedPayload.Username)
	require.NotZero(t, returnedPayload.ID)
	require.WithinDuration(t, issuedAt, returnedPayload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, returnedPayload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomUsername()
	duration := -time.Minute

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	returnedPayload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, err, ErrExpiredToken)
	require.Nil(t, returnedPayload)
}

func TestInvalidToken(t *testing.T) {
	username := util.RandomUsername()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	require.NotEmpty(t, jwtToken)

	invalidToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, invalidToken)

	maker, err := NewJWTMaker(util.RandomString(32, util.ALL))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.Equal(t, ErrInvalidToken, err)
	require.Nil(t, payload)
}
