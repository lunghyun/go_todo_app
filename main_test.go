package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

//- 예상한 대로 HTTP 서버가 실행되는가?
//- 테스트 코드가 의도한 대로 종료 처리를 하는가?

//코드의 흐름은 아래와 같다.
//- 취소 가능한 ctx 객체를 만든다
//- 다른 고루틴에서 테스트 대상인 run을 실행, HTTP 서버를 시작
//- 엔드포인트에 대해 GET 요청 전송
//- cancel 함수 실행
//- *errgroup.Group.Wait 메서드의 경유로 run 함수의 반환값을 검증한다.
//- Get 요청에서 받은 응답 바디가 기대한 문자열인 것을 검증한다

func TestRun(t *testing.T) {
	//- 취소 가능한 ctx 객체를 만든다
	ctx, cancel := context.WithCancel(context.Background())

	//- 다른 고루틴에서 테스트 대상인 run을 실행, HTTP 서버를 시작
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	//- 엔드포인트에 대해 GET 요청 전송
	in := "message"
	rsp, err := http.Get("http://localhost:8080/" + in)
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
