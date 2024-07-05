package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
)

// HLS 스트림을 수신하여 저장하는 핸들러
func UploadHLSHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(r) // 요청 정보를 로그에 기록

	if r.Method == http.MethodPut {
		// HLS 스트림 파일 경로 설정
		filePath := "." + r.URL.Path
		log.Println("Receiving stream:", filePath)

		file, err := os.Create(filePath)
		if err != nil {
			log.Println("Failed to create file:", err)
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}

		// defer를 사용하여 파일 닫기, 닫는 동안 발생하는 오류 처리
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println("Failed to close file:", err)
			}
		}(file)

		// 요청 바디 읽어서 파일에 쓰기
		_, err = io.Copy(file, r.Body)
		if err != nil {
			log.Println("Failed to write to file:", err)
			http.Error(w, "Failed to write to file", http.StatusInternalServerError)
			return
		}

		log.Println("Stream saved:", filePath)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
