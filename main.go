package main

import (
	"github.com/HadryanSilva/go-youtube-downloader/core/ui"
	"github.com/HadryanSilva/go-youtube-downloader/core/utils"
)

func main() {
	utils.InstallDependencies()
	ui.GenerateWindow()
}
