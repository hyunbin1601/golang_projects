package main

// ioutil 패키지를 이용한 파일 i/o

import "io/ioutil" // ioutil 패키지는 더이상 사용되지 않음
// 쓸 수는 있으나 권장되지 않으므로 os 패키지의 os.ReadFile, os.WriteFile 사용하기

func Read_and_write_ioutil() {
	bytes, err := ioutil.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("test2.txt", bytes, 0) // 0은 파일 권한, 기본값으로 설정, 어떤 권한? -> 파일 읽고 쓰는 권한
	if err != nil {
		panic(err)
	}
}