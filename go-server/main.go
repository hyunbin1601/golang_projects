package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) { // *http.Request는 http 요청을 나타내는 구조체, 여기서 r을 통해 요청에서 데이터를 읽을 수 있음
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	name := r.FormValue("name") // name 변수에 name 키의 값을 저장
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "address = %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) { // *http.Request는 http 요청을 나타내는 구조체
	// fmt.Fprintf(w, "Hello, World!")
	if r.URL.Path != "/hell" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello, World!")

}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler) // formHandler 함수를 /form 경로에 등록
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	} // if 문 이후의 세미콜론은 변수 선언 후 조건을 평가할 때 둘을 구분하기 위해 사용
	// 이해가 안가겠지만...세미콜론 뒤에 오는 부분이 실제 if 조건이다
	// 해석 -> err := http.ListenAndServe(":8080", nil) 문으로 err에 서버 시작시 만약 에러 발생하면 변수 err에 저장, 그리고 세미콜론 후 실제 if 조건을 넣어 err를 평가함
}
