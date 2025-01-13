package utils

import (
	"os"
	"os/exec"
	"strings"
)

func InstallDependencies() {
	env := os.Getenv("PATH")
	values := strings.Split(env, ";")
	for _, value := range values {
		if strings.Contains(value, "Gyan.FFmpeg_Microsoft.Winget.Source") &&
			strings.Contains(value, "yt-dlp.yt-dlp_Microsoft.Winget.Source") {
			return
		}
	}

	exec.Command("winget", "install", "yt-dlp")
}
