package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

const (
	maxSize int = 1024 * 1024 * 50
)

func getFileContent(fileName string) ([]byte, int, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return []byte{}, 0, errors.New("file not found")
	}

	defer fd.Close()

	reader := bufio.NewReader(fd)
	info, _ := fd.Stat()
	data := make([]byte, info.Size()) // BUG: get correct size here
	count, err := io.ReadFull(reader, data)
	if err != nil {
		return []byte{}, 0, errors.New("file read failed")
	}

	log.Printf("bytes read = %d", count)
	return data, count, nil
}

func streamManifestHandler(ctx *fiber.Ctx) {
	log.Println("manifest request received")
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctx.SendStatus(404)
		return
	}

	fileName := fmt.Sprintf("public/%d/index.m3u8", id)
	data, _, err := getFileContent(fileName)
	if err != nil {
		ctx.SendStatus(404)
		return
	}

	ctx.Set("Content-Type", "application/x-mpegurl")
	ctx.Write(data)
}

func streamSegmentHandler(ctx *fiber.Ctx) {
	log.Println("segment request received")
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctx.SendStatus(404)
		return
	}

	segment := ctx.Params("segment")

	fileName := fmt.Sprintf("public/%d/%s", id, segment)
	data, _, err := getFileContent(fileName)
	if err != nil {
		ctx.SendStatus(404)
		return
	}

	ctx.Set("Content-Type", "video/mp2t")
	ctx.Write(data)
}

// StartServer starts the media streaming server.
func StartServer() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/media", func(c *fiber.Ctx) {
		c.Send("Hello World")
	})

	app.Get("/media/:id/stream/", streamManifestHandler)
	app.Get("/media/:id/:segment", streamSegmentHandler)

	app.Static("/", "./public")
	app.Listen(8080)
}
