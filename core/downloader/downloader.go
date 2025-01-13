package downloader

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func DownloadVideo(info DownloadInfo, progressChan chan float64) error {
	valid, err := validateResolution(info.Url, info.Resolution)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("resolution %sp not found", info.Resolution)
	}

	cmd := exec.Command(
		"yt-dlp",
		"-S", "res:"+info.Resolution,
		"-f", "bv+ba/b",
		"-o", info.Path+"/%(title)s.%(ext)s",
		info.Url,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not create output pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start yt-dlp: %v", err)
	}

	progressRegexp := regexp.MustCompile(`(\d+)%`)

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("could not read output of yt-dlp: %v", err)
		}

		output := string(buf[:n])
		if match := progressRegexp.FindStringSubmatch(output); match != nil {
			progress := match[1]
			var p float64
			fmt.Sscanf(progress, "%f", &p)
			progressChan <- p
		}

		if err != nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("could not wait command: %v", err)
	}

	return nil
}

func validateResolution(url string, resolution string) (bool, error) {
	cmd := exec.Command("yt-dlp", "--list-formats", url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("could not create output pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("could not start yt-dlp: %v", err)
	}

	buf := make([]byte, 1024)
	var availableResolutions []string
	for {
		n, err := stdout.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return false, fmt.Errorf("could not read output of yt-dlp: %v", err)
		}

		output := string(buf[:n])

		resolutionRegexp := regexp.MustCompile(`(\d+p)`)
		matches := resolutionRegexp.FindAllString(output, -1)
		for _, match := range matches {
			availableResolutions = append(availableResolutions, match)
		}

		if err != nil {
			break
		}
	}

	for _, res := range availableResolutions {
		if strings.Contains(res, resolution) {
			return true, nil
		}
	}

	return false, fmt.Errorf("the resolution %sp are not available for this video", resolution)
}
