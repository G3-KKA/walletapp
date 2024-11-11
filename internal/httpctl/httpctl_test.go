package httpctl

import (
	"context"
	"net"
	"testing"
	"time"
	"walletapp/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// Gin design is insane, you either do NON-Parallel tests, or you see those,
// not that handsome messages from gin-debug, that solution is thread-safe and works.
//
// Hint: it should be put in every test that uses gin.
//
// It SHOULDNT affect production code because test functions are not compiled into
// any build of actual program.
var _ = func() struct{} { gin.SetMode(gin.ReleaseMode); return struct{}{} }()

const (
	// httpctl_test actually uses port, might come with the problems.
	//
	// Yes, port is guaranteed to be free-to-use, but be careful.
	_IngorePotentiallyDangerousTest = false

	// This test shouldn't (?) be run in parallel,
	// so any not-immediate value of [_TimeSleep] will globally affect testing time.
	_TimeSleep = time.Second * 1
)

func TestHTTPCtl(t *testing.T) {

	if _IngorePotentiallyDangerousTest {
		t.SkipNow()
	}

	l := zerolog.New(zerolog.NewTestWriter(t))

	getFreePortString := func() string {
		t.Helper()

		a, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			t.Skip("was not be able to get any free port from os, err:", err)

			return ""
		}

		l, err := net.ListenTCP("tcp", a)
		if err != nil {
			t.Skip("was not be able to get any free port from os, err:", err)

			return ""
		}

		defer l.Close()

		return l.Addr().(*net.TCPAddr).AddrPort().String()
	}

	cfg := config.HTTPServer{
		Address: getFreePortString(),
	}

	httpC := New(&l, cfg, gin.New())

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(_TimeSleep)
		cancel()
	}()

	err := httpC.Serve(ctx)

	assert.NoError(t, err)

}
