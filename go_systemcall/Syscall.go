package main

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func Syscall() {
	// window의 패키지 이용
	hWnd := uintptr(0)
	windows.MessageBox(
		windows.HWND(hWnd), // [in, optional] HWND    hWnd,
		windows.StringToUTF16Ptr("Used Windows Package"), // [in, optional] LPCTSTR lpText,
		windows.StringToUTF16Ptr("MessageBox 1/2"),       // [in, optional] LPCTSTR lpCaption,
		windows.MB_OK) // [in] UINT    uType

	// syscall 패키지 이용
	user32dll := syscall.NewLazyDLL("User32.dll")
	procMsgBox := user32dll.NewProc("MessageBoxW")  // 프로시저(함수)를 가져옴

	hWnd = uintptr(0)
	lpText, err := syscall.UTF16PtrFromString("Used Syscall Package")
	if err != nil {
		log.Fatalln("lpText UTF16PtrFromString Failed")
	}
	lpCaption, err := syscall.UTF16PtrFromString("MessageBox 2/2")
	if err != nil {
		log.Fatalln("lpCaption UTF16PtrFromString Failed")
	}
	uType := uint(0)

	procMsgBox.Call(
		hWnd, 
		uintptr(unsafe.Pointer(lpText)), 
		uintptr(unsafe.Pointer(lpCaption)), 
		uintptr(uType))
}