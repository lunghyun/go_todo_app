package main

import (
	"fmt"
	"net/http"
	"os"
)

// 요청을 받아 응답 메시지를 생성하는 서버
// 포트번호 8080 고정

func main() {

	// ListenAndServe:
	// arg1: 주소 문자열(IP 정보 생략 -> local)
	// arg2: 핸들러(단일 처리만 구현 -> 어떤 path든 path를 통한 응답 메시지 반환
	err := http.ListenAndServe(":8080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	)
	// 기본적으로 병렬 요청 처리됨
	if err != nil {
		fmt.Printf("failed to start server: %v", err)
		os.Exit(1)
	}
}
