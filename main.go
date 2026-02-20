package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("need Port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", p, err)
	}
	if err = run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 인수로 받은 net.Listener를 이용하므로 Addr 필드는 지정하지 않는다.
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ListenAndServe -> Serve로 변경
		// http.ErrServerClosed는 httpServerShutdown이 정상 종료되었다고 표시하므로 문제없음
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
