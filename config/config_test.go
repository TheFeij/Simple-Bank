package config

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./", "config_test")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	require.Equal(t, "http host", config.HTTPServerHost)
	require.Equal(t, "http port", config.HTTPServerPort)
	require.Equal(t, "grpc host", config.GrpcServerHost)
	require.Equal(t, "grpc port", config.GrpcServerPort)
	require.Equal(t, "source", config.DatabaseSource)
	require.Equal(t, "driver", config.DatabaseDriver)
	require.Equal(t, "key", config.TokenSymmetricKey)
	require.Equal(t, "environment", config.Environment)
	require.Equal(t, 1*time.Minute, config.TokenAccessTokenDuration)
}
