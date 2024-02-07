package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	exitCode := m.Run()
	os.Exit(exitCode)
}
