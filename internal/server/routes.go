package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"strings"
	app "yourtube/internal/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.healthHandler)

	fmt.Println("checking RegisterRoutes")
	fmt.Println(e.POST("/upload", s.handelUpload))

	e.Start(":8080")

	return e
}

func isValidFileType(filename string, allowedExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func (s *Server) handelUpload(c echo.Context) error {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File not found in the request"})
	}
	// Specify the absolute path to the 'internal' directory
	internalDir := "./internal"
	// Create the destination directory if it doesn't exist
	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		err := os.Mkdir(internalDir, 0755)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create destination directory"})
		}
	}
	// Now create the 'input' directory inside 'internal'
	dirname := filepath.Join(internalDir, "input")
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create destination directory"})
		}
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check destination directory"})
	} else {
		log.Printf("Directory %s already exists", dirname)
	}
	// Use filepath.Join to properly concatenate directory and filename
	dstPath := filepath.Join(dirname, file.Filename)

	// Get the absolute path
	absDstPath, err := filepath.Abs(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get absolute path"})
	}

	dst, err := os.Create(absDstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create destination file"})
	}
	defer dst.Close()

	// Open the source file from the request
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open source file"})
	}
	defer src.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to copy file contents"})
	}

	app.Transcoder(absDstPath)
	// Handle further processing, if needed
	// fmt.Println(dstPath)
	return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully"})
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
