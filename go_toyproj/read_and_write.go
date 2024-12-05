package main

import (
	"io"
	"os"
)

// 파일 읽고 쓰기
// os 패키지 이용

func Read_and_write() {
	fileInput, err := os.Open("test.txt")
	if err != nil {
		panic(err) // panic은 에러가 발생하면 프로그램을 중단시키는 함수
	}
	defer fileInput.Close()

	fileOutput, err := os.Create("test2.txt")

	if err != nil {
		panic(err)
	}
	defer fileOutput.Close()

	buffer := make([]byte, 1024) // 1024바이트씩 읽기

	for {
		cnt, err := fileInput.Read(buffer) // 파일 읽기, buffer에 읽은 바이트 수 저장
		if err != nil && err != io.EOF {
			panic(err)
		}
		if cnt == 0 {
			break // 루프 종료, 다 읽었다는 의미
		}

		_, err = fileOutput.Write(buffer[:cnt]) // 파일 쓰기, _는 사용하지 않는 리턴값
		if err != nil {
			panic(err)
		}
	}
}




