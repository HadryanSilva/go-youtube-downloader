package downloader

import (
	"fmt"
	"os/exec"
	"regexp"
)

func DownloadVideo(url string, path string, progressChan chan float64) error {
	cmd := exec.Command("yt-dlp", "-S", "res:720", "-f", "bv+ba/b", "-o", path+"/%(title)s.%(ext)s", url)

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
