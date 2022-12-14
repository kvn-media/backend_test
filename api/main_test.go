package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/kvn-media/backend_test/db/sqlc"
	"github.com/kvn-media/backend_test/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// set gin into TestMode to get cleaner logs
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
