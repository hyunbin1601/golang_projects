package main

import (
	"fmt"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows"
)

// APC 콜백 함수
func APCFunction(par uintptr) uintptr {
	fmt.Printf("APC 실행중 / 전달된 값: %d\n", par)
	return 0
}

// QueueUserAPC 호출을 위한 준비
var modKernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procQueueUserAPC = modKernel32.NewProc("QueueUserAPC")

// QueueUserAPC 함수 래핑
func QueueUserAPC(pfnAPC uintptr, hThread uintptr, dwData uintptr) error {
	ret, _, err := procQueueUserAPC.Call(pfnAPC, hThread, dwData)
	if ret == 0 {
		return err
	}
	return nil
}

func main() {
	// OS 스레드 고정
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// 현재 스레드 핸들 가져오기
	currentThread, err := windows.GetCurrentThread()
	if err != nil {
		fmt.Printf("현재 스레드 가져오기 실패: %v\n", err)
		return
	}

	// APC 등록
	fmt.Println("스레드에 APC 등록")
	for i := uintptr(1); i <= 3; i++ {
		err := QueueUserAPC(
			syscall.NewCallback(APCFunction), // APC 콜백 함수
			uintptr(currentThread),           // 타겟 스레드 핸들
			i,                                // 전달할 값
		)
		if err != nil {
			fmt.Printf("APC 등록 실패: %v\n", err)
			return
		}
	}

	// Alertable 상태로 진입
	fmt.Println("Alertable 상태로 진입")
	for i := 0; i < 5; i++ {
		fmt.Printf("대기 중 (%d)...\n", i+1)
		windows.SleepEx(1000, true) // Alertable 상태로 대기
	}

	fmt.Println("프로그램 종료")
}
