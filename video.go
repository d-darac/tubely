package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getVideoAspectRatio(filePath string) (string, error) {
	type FFProbe struct {
		Streams []struct {
			Width  int64 `json:"width"`
			Height int64 `json:"height"`
		} `json:"streams"`
	}

	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-print_format",
		"json",
		"-show_streams",
		filePath,
	)

	buf := bytes.Buffer{}
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	ffprobe := FFProbe{}
	if err := json.Unmarshal(buf.Bytes(), &ffprobe); err != nil {
		return "", err
	}

	divisor := gcd(ffprobe.Streams[0].Height, ffprobe.Streams[0].Width)
	return fmt.Sprintf(
		"%d:%d",
		(ffprobe.Streams[0].Height / divisor),
		(ffprobe.Streams[0].Width / divisor),
	), nil

}

func getVideoOrientation(aspectRatio string) (string, error) {
	split := strings.Split(aspectRatio, ":")
	if len(split) < 2 {
		return "", fmt.Errorf("invalid aspect ratio")
	}

	height, err := strconv.Atoi(split[0])
	if err != nil {
		return "", fmt.Errorf("invalid aspect ratio")
	}
	width, err := strconv.Atoi(split[1])
	if err != nil {
		return "", fmt.Errorf("invalid aspect ratio")
	}

	if height > width {
		return "portrait", nil
	} else if width > height {
		return "landscape", nil
	} else {
		return "other", nil
	}
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
