package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen on port %v", err)
	}
	//- 취소 가능한 ctx 객체를 만든다
	ctx, cancel := context.WithCancel(context.Background())

	//- 다른 고루틴에서 테스트 대상인 run을 실행, HTTP 서버를 시작
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	//- 엔드포인트에 대해 GET 요청 전송
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to GET: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP 서버 반환값 검증
	want := fmt.Sprintf("Hello, %s", in)
	if string(got) != want {
		t.Errorf("want %s, but got %s", want, got)
	}
	//- cancel 함수 실행
	cancel()

	// run 함수 반환값 검증
	if err = eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
