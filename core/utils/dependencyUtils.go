package utils

import (
	"os"
	"os/exec"
	"strings"
)

func InstallDependencies() {
	if os.Getenv("os") == "windows" {
		env := os.Getenv("PATH")
		values := strings.Split(env, ";")
		for _, value := range values {
			if strings.Contains(value, "FFmpeg_Microsoft.Winget.Source") &&
				strings.Contains(value, "yt-dlp.yt-dlp_Microsoft.Winget.Source") {
				return
			}
		}
		exec.Command("winget", "install", "yt-dlp")
	} else if os.Getenv("os") == "linux" {
		_, err := exec.Command("which", "yt-dlp").Output()
		if err != nil {
			exec.Command("sudo", "add-apt-repository", "ppa:tomtomtom/yt-dlp")
			exec.Command("sudo", "apt", "update")
			exec.Command("sudo", "apt", "install", "yt-dlp")
		}
	}
}
