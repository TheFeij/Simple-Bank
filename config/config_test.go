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

	require.Equal(t, "host", config.ServerHost)
	require.Equal(t, "port", config.ServerPort)
	require.Equal(t, "source", config.DatabaseSource)
	require.Equal(t, "driver", config.DatabaseDriver)
	require.Equal(t, "key", config.TokenSymmetricKey)
	require.Equal(t, "development", config.Environment)
	require.Equal(t, 1*time.Minute, config.TokenAccessTokenDuration)
}
