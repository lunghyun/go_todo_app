package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ListenAndServe -> Serve로 변경
		// http.ErrServerClosed는 httpServerShutdown이 정상 종료되었다고 표시하므로 문제없음
		if err := s.srv.Serve(s.l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// ch로부터 종료 알림을 기다림
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Go 메서드로 실행한 다른 goroutine의 종료를 기다림
	return eg.Wait()
}
