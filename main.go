package main

import (
	"SnackCam/handlers"
	"log"
	"net/http"
)

func main() {
	// 메세지 전달용 웹소켓 실행
	go handlers.HandleMessages()

	// 웹 소켓 핸들러 설정
	http.HandleFunc("/ws", handlers.HandleConnections)

	// 서버 시작 시 Gstreamer 실행
	if err := handlers.StartGstreamer(); err != nil {
		log.Fatalf("Failed to start Gstreamer: %v", err)
	}

	// HLS 파일이 있는 디렉토리를 설정합니다.
	http.HandleFunc("/hls/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		fs := http.StripPrefix("/hls", http.FileServer(http.Dir("static/hls")))
		fs.ServeHTTP(w, r)
	})

	// HTML 파일이 있는 디렉토리를 설정하고, 로그를 추가합니다.
	http.HandleFunc("/", handlers.FileServerHandler)

	// HLS 스트림을 수신하여 저장하는 핸들러 설정
	http.HandleFunc("/upload/hls/", handlers.UploadHLSHandler)

	log.Println("Starting server on :18080")
	if err := http.ListenAndServe(":18080", nil); err != nil {
		log.Fatal(err)
	}
}
