package main

import (
	"golang.org/x/sys/windows"
)

func WindowApi() {
    hWnd := uintptr(0)
    windows.MessageBox(
        windows.HWND(hWnd), // Handle to the owner window
        windows.StringToUTF16Ptr("Using window package"),
        windows.StringToUTF16Ptr("MessagBox Hi"),
        windows.MB_OK,  // Flags -> OK button
    )
}
