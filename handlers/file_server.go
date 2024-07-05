package handlers

import (
	"net/http"
	"path/filepath"
)

// 파일 서버 요청을 처리하는 핸들러
func FileServerHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(r) // 요청 정보를 로그에 기록

	// 파일 서버 기능 수행
	// 요청된 파일 경로 설정
	requestedPath := filepath.Join("static/html", r.URL.Path)

	//// 경로가 HTML 파일로 끝나는지 확인
	//if !strings.HasSuffix(requestedPath, ".html") {
	//	http.Error(w, "Forbidden", http.StatusForbidden)
	//	return
	//}

	http.ServeFile(w, r, requestedPath)
}
