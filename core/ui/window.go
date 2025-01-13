package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/HadryanSilva/go-youtube-downloader/core/downloader"
	"strings"
)

func GenerateWindow() {
	myApp := app.New()
	window := myApp.NewWindow("YouTube Video Downloader")
	window.Resize(fyne.NewSize(820, 640))

	iconResource, err := fyne.LoadResourceFromPath("icon.ico")
	if err != nil {
		fmt.Errorf("erro ao carregar o ícone: %v", err)
	}
	window.SetIcon(iconResource)

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Cole o link do YouTube aqui")

	var resolutionSelected string
	selectResolutionMenu := widget.NewSelect([]string{"480p", "720", "1080p", "1440p", "2160p"}, func(resolution string) {
		resolutionSelected = strings.TrimSuffix(resolution, "p")
	})
	selectResolutionMenu.PlaceHolder = "Selecione a resolução"

	var outputPath string
	pathLabel := widget.NewLabel("Pasta de destino: Não selecionada")

	selectFolderButton := widget.NewButton("Selecionar pasta", func() {
		newFolderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri == nil {
				return
			}
			outputPath = uri.Path()
			pathLabel.SetText("Pasta de destino: " + outputPath)
		}, window)
		newFolderDialog.Show()
	})
	selectFolderButton.Resize(fyne.NewSize(100, 100))

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	downloadButton := widget.NewButton("Baixar Vídeo", func() {
		if urlEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("por favor, insira uma URL válida"), window)
			return
		}
		if outputPath == "" {
			dialog.ShowError(fmt.Errorf("por favor, selecione uma pasta de destino"), window)
			return
		}
		if resolutionSelected == "" {
			dialog.ShowError(fmt.Errorf("por favor, selecione uma resolução"), window)
			return
		}

		progressBar.Show()
		progressChan := make(chan float64)
		info := downloader.DownloadInfo{
			Url:        urlEntry.Text,
			Path:       outputPath,
			Resolution: resolutionSelected,
		}
		go downloadVideo(info, window, progressBar, progressChan)
	})

	content := container.NewVBox(
		widget.NewLabel("Baixar Vídeo do YouTube"),
		urlEntry,
		selectResolutionMenu,
		selectFolderButton,
		pathLabel,
		downloadButton,
		progressBar,
	)

	window.SetContent(content)
	window.ShowAndRun()
}

func downloadVideo(info downloader.DownloadInfo, window fyne.Window, progress *widget.ProgressBar, progressChan chan float64) {
	defer close(progressChan)

	go func() {
		for p := range progressChan {
			progress.SetValue(p / 100)
		}
	}()

	err := downloader.DownloadVideo(info, progressChan)

	progress.Hide()
	if err != nil {
		dialog.ShowError(fmt.Errorf("erro ao baixar vídeo: %v", err), window)
		return
	}
	dialog.ShowInformation("Sucesso", "Download concluído com sucesso!", window)
}
