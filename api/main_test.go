package api

import (
	"Simple-Bank/config"
	"Simple-Bank/db/services"
	"Simple-Bank/token"
	"Simple-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var configs *config.Config

func NewTestServer(t *testing.T, services services.Services, tokenMaker token.Maker) *Server {
	server, err := NewServer(getTestConfig(), services, tokenMaker)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func getTestConfig() *config.Config {
	return &config.Config{
		TokenAccessTokenDuration:  15 * time.Minute,
		TokenRefreshTokenDuration: 24 * time.Hour,
		TokenSymmetricKey:         util.RandomString(32, util.ALL),
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	configs = getTestConfig()

	exitCode := m.Run()
	os.Exit(exitCode)
}
