package handlers

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// UploadHandler handles file uploads and converts them to HLS format
var mux sync.Mutex

// UploadHandler handles file uploads and converts them to HLS format
func UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to retrieve file from form-data")
	}

	// Save the uploaded file to the server
	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create temporary file")
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open uploaded file")
	}
	defer uploadedFile.Close()

	if _, err := io.Copy(tempFile, uploadedFile); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	mux.Lock()
	defer mux.Unlock()

	// Convert the uploaded file to HLS format and append to existing stream
	hlsDir := "static/hls"
	absHlsDir, err := filepath.Abs(hlsDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get absolute path of HLS directory")
	}

	segmentFile := filepath.Join(absHlsDir, "temp_segment.m3u8")
	cmd := exec.Command("ffmpeg", "-i", tempFilePath, "-codec: copy", "-start_number", "0",
		"-hls_time", "10", "-hls_list_size", "0", "-f", "hls", segmentFile)

	if err := cmd.Run(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to convert video to HLS format")
	}

	// Append new segments to the main playlist
	if err := appendToPlaylist(absHlsDir, segmentFile); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update HLS playlist")
	}

	return c.JSON(fiber.Map{"status": "success"})
}

// appendToPlaylist appends segments from the new file to the existing playlist
func appendToPlaylist(hlsDir, segmentFile string) error {
	mainPlaylist := filepath.Join(hlsDir, "playlist.m3u8")
	//tempPlaylist := mainPlaylist + ".tmp"

	// Open main playlist for appending
	main, err := os.OpenFile(mainPlaylist, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer main.Close()

	// Read the new segment file and append to main playlist
	segment, err := os.Open(segmentFile)
	if err != nil {
		return err
	}
	defer segment.Close()

	// Read the new segment file
	newSegments, err := io.ReadAll(segment)
	if err != nil {
		return err
	}

	// Write the new segments to the main playlist
	if _, err := main.Write(newSegments); err != nil {
		return err
	}

	return nil
}
