package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows" // windows 패키지는 Windows API를 사용할 수 있게 해줌
)

var startupInfo windows.StartupInfo
var outProcInfo windows.ProcessInformation

func Create_Process() {
	path := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe" // path 변수 선언 및 초기화(golang 컴파일러가 타입 추론)

	err := windows.CreateProcess(
		nil,
		windows.StringToUTF16Ptr(path),
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&startupInfo,
		&outProcInfo,
	)

	if err != nil {
		log.Fatalf("CreateProcess failed: %v", err) // %v는 값의 기본 형식을 사용하여 값을 출력, %v는 구체적인 형식을 지정하지 않아도 됨
	}

	fmt.Printf("Process Created from path: %s PID: %d\n", path, outProcInfo.ProcessId)
	fmt.Printf("Process Handle: %x \n Thread Handle: %x\n", outProcInfo.Process, outProcInfo.Thread)
}