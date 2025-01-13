## Golang Youtube Video Downloader UI
This is a simple YouTube video downloader written in Golang and was made for a complete beginner in this language. 
It uses the yt-dlp library to download the videos.

### Requirements
- C compiler (gcc) for the fyne library

### Generate Executable
To generate the executable, you can use the following command:
#### For Windows
```bash
  fyne package -os windows -name youtube-downloader -icon icon.ico
```
#### For Linux
```bash
  fyne package -os linux -name youtube-downloader -icon icon.ico
```

