package api

import (
	"Simple-Bank/config"
	"Simple-Bank/db/services"
	"Simple-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func NewTestServer(t *testing.T, services services.Services) *Server {
	config := &config.Config{
		TokenAccessTokenDuration: time.Minute,
		TokenSymmetricKey:        util.RandomString(32, util.ALL),
	}

	server, err := NewServer(config, services)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	exitCode := m.Run()
	os.Exit(exitCode)
}
