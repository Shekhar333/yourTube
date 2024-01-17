package controllers

import (
	// "bufio"
	"fmt"
	// "os"
	"os/exec"
	"sync"
	"time"
)

func Execute(dst string, i int, pScan int, wg *sync.WaitGroup) {
	defer wg.Done()

	outputPath := fmt.Sprintf("./output/output%d.mp4", i)
	// absDstPath, err := filepath.Abs(dstPath)
	// if err != nil {
	// 	return JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get absolute path"})
	// }
	scaleFilter := fmt.Sprintf("scale=trunc(oh*a/2)*2:%d", pScan)
	cmd := exec.Command("ffmpeg", "-i", dst, "-vf", scaleFilter, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-strict", "experimental", outputPath)

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error:", i, err)
		return
	}
}

func Transcoder(input string) {
	progressiveScan := [6]int{1080, 720, 480, 360, 240, 144}
	startTime := time.Now()
	var wg sync.WaitGroup

	// line := scanner.Text()
	for i, pScan := range progressiveScan {
		wg.Add(1)
		go Execute(input, i, pScan, &wg)
	}

	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Command executed successfully in: ", elapsedTime)
}
