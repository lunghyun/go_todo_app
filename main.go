package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/lunghyun/go_todo_app/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

// run starts the server: it loads configuration, opens a TCP listener, builds a request mux, and runs the server.
//
// If configuration creation fails, the error is returned. If the listener cannot be created the process exits
// after logging a fatal message. Resources produced by NewMux are released via the deferred cleanup even if
// NewMux returns an error; however, a NewMux error is not propagated. Any error returned by the server's Run
// is returned to the caller.
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux, cleanup, err := NewMux(ctx, cfg)
	// 오류 반환되어도 cleanup
	defer cleanup()
	if err != nil {
		_ = fmt.Errorf("NewMux failed: %v", err)
	}
	s := NewServer(l, mux)
	return s.Run(ctx)
}
