package main

import (
	"SnackCam/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})
	app.Use(cors.New())
	app.Use(logger.New())

	// 메세지 전달용 웹소켓 실행
	go handlers.HandleMessages()

	// 웹 소켓 핸들러 설정
	app.Get("/ws", websocket.New(handlers.HandleConnections))

	/////////////////////////////////////////////////////// 카메라에서 다이렉트로 전송 받는 경우

	// 서버 시작 시 Gstreamer 실행
	//if err := handlers.StartGstreamer(); err != nil {
	//	log.Fatalf("Failed to start Gstreamer: %v", err)
	//}

	//// HLS 스트림을 수신하여 저장하는 핸들러 설정
	//app.Post("/upload/hls", func(c *fiber.Ctx) error {
	//	return handlers.UploadHLSHandler(c)
	//})

	///////////////////////////////////////////////////////

	// HLS 파일이 있는 디렉토리를 설정합니다.
	app.Use("/hls", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Cache-Control", "no-cache")
		return c.Next()
	})
	app.Static("/hls", "static/hls")

	// HTML 파일이 있는 디렉토리를 설정하고, 로그를 추가합니다.
	app.Get("/", func(c *fiber.Ctx) error {
		return handlers.FileServerHandler(c)
	})

	// 비디오 업로드 -> HLS 변환
	app.Post("/uploadVideo", handlers.UploadHandler)

	log.Println("Starting server on :18080")
	if err := app.Listen(":18080"); err != nil {
		log.Fatal(err)
	}

}
