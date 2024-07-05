package handlers

import (
	"log"
	"os/exec"
	"sync"
)

const (
	rtpPort      = "15000"
	hlsDirectory = "<절대경로>/SnackCam/static/hls" // TODO --- 절대 경로 부분 동적으로 추가해줘야 함
)

var (
	cmd    *exec.Cmd
	cmdMux sync.Mutex
)

func StartGstreamer() error {
	cmdMux.Lock()
	defer cmdMux.Unlock()

	// 이미 실행 중인 Gstreamer 프로세스가 있으면 종료
	if cmd != nil {
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
	}

	// Gstreamer 명령어 구성
	gstCmd := exec.Command("gst-launch-1.0",
		"-e", "udpsrc", "port="+rtpPort, "caps=application/x-rtp, payload=96",
		"!", "queue", "max-size-buffers=1000", "max-size-time=0", "max-size-bytes=0",
		"!", "rtph264depay",
		"!", "h264parse",
		"!", "mpegtsmux",
		"!", "hlssink", "location="+hlsDirectory+"/segment%05d.ts", "playlist-location="+hlsDirectory+"/playlist.m3u8", "target-duration=5", "max-files=5", "playlist-length=5",
	)

	// 표준 출력과 오류를 연결
	gstCmd.Stdout = log.Writer()
	gstCmd.Stderr = log.Writer()

	// Gstreamer 명령어 실행
	if err := gstCmd.Start(); err != nil {
		return err
	}

	cmd = gstCmd
	return nil
}
