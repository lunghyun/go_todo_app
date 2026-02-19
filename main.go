package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	// 다른 고루틴에서 HTTP 서버를 실행
	eg.Go(func() error {
		// http.ErrServerClosed:
		// Shutdown이 정상 종료된 것을 나타내므로 이상 처리가 아님
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// ch로부터 종료 알림을 기다림
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Go 메서드로 실행한 다른 goroutine의 종료를 기다림
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
