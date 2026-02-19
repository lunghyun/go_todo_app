package main

import "testing"

func TestMainFunc(t *testing.T) {
	go main()
	// 실행은 가능, 종료 지시 불가능
	// 문제점
	// 1. 테스트 완료 후 종료 방법 없음
	// 2. 출력 검증 어려움
	// 3. 이상 처리 시, os.Exit 함수 호출됨
	// 4. 포트번호 고정 -> 테스트 서버 실행이 안될 수 있음
}
