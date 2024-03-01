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

	require.Equal(t, "host", config.Server.Host)
	require.Equal(t, "port", config.Server.Port)
	require.Equal(t, "source", config.Database.Source)
	require.Equal(t, "driver", config.Database.Driver)
	require.Equal(t, "key", config.Token.SymmetricKey)
	require.Equal(t, 1*time.Minute, config.Token.AccessTokenDuration)
}
