package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

// HLS 스트림을 수신하여 저장하는 핸들러
func UploadHLSHandler(c *fiber.Ctx) error {
	// 파일 처리 로직을 추가합니다.
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	log.Println("Received file:", file.Filename)
	// 파일 저장 등의 로직을 추가합니다.

	return c.SendString("File uploaded successfully")
}
