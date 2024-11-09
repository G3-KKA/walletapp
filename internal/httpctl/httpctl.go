package httpctl

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"walletapp/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type httpController struct {
	server *http.Server
	l      *zerolog.Logger
	cfg    config.HTTPServer

	shutdown sync.Once
}

// New is a httpController constructor.
func New(l *zerolog.Logger, cfg config.HTTPServer, mux *gin.Engine) *httpController {
	server := &http.Server{
		Addr:                         cfg.Address,
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     log.New(l, "", log.Flags()),
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	httpC := &httpController{
		server:   server,
		l:        l,
		cfg:      cfg,
		shutdown: sync.Once{},
	}

	return httpC
}

// Serve starts listening.
// May be shutted down via context.
func (ctl *httpController) Serve(ctx context.Context) error {

	lis, err := net.Listen("tcp", ctl.cfg.Address)
	if err != nil {
		return err
	}

	// Straight and simple way to accumulate all errors, no channels required.
	errmx := sync.Mutex{}
	accumulateError := func(target error) {
		if !errors.Is(target, http.ErrServerClosed) {
			errmx.Lock()
			err = errors.Join(target, err)
			errmx.Unlock()
		}
	}

	ctl.l.Info().Msg("http server starting on: " + lis.Addr().String())

	servRoutineExited := make(chan struct{})
	serverRoutine := func() {

		servErr := ctl.server.Serve(lis)
		accumulateError(servErr)

		close(servRoutineExited)
	}
	go serverRoutine()

	select {
	case <-servRoutineExited:
	case <-ctx.Done():

		shuterr := ctl.Shutdown(context.TODO())
		accumulateError(shuterr)

		// Still waiting for server itself to close.
		<-servRoutineExited
	}

	// ctl.Shutdown() is safe to call as much times as needed.
	shuterr := ctl.Shutdown(context.TODO())
	accumulateError(shuterr)

	if err != nil {

		return err
	}

	return nil

}

/* 	ctxshutdowner := func() {
   		select {
   		case <-ctx.Done():
   			errs <- ctl.Shutdown(context.TODO())
   		case <-doneTrigger:
   			// Eventually it goroutine will be terminated,
   			// even if context cannot be canceled.
   		}
   		close(shutdownerExited)
   	}
   	go ctxshutdowner() */

// Shutdown httpController gracefully.
func (ctl *httpController) Shutdown(ctx context.Context) error {

	var err error

	ctl.shutdown.Do(func() {
		err = ctl.server.Shutdown(ctx)
	})

	if err != nil {
		return err
	}

	return nil
}
