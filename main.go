package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func run(ctx context.Context) error {
	// HTTP 서버 실행 시, *http.Server 타입을 경유하는 것이 정석
	// - Shutdown 메서드 존재
	// - ListenAndServe 존재
	// - ctx로 외부에서 취소 처리를 받을 수 있음
	s := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}
	_ = s.ListenAndServe()

	return nil
}

func main() {
	if err := run(context.TODO()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
