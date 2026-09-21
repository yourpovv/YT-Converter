package main

import (
	"YTConverter/internal/banner"
	"YTConverter/internal/download"
	"YTConverter/utils"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type downloader func(youtubeURL, outputPath string) error

type downloadOption struct {
	extension     string
	progressLabel string
	fetch         downloader
}

func main() {
	banner.SetTitle("YT Converter")
	runMenu()
}

func runMenu() {
	mp3 := downloadOption{extension: download.ExtMP3, progressLabel: "[*] Downloading and converting to MP3...", fetch: download.DownloadMP3}
	mp4 := downloadOption{extension: download.ExtMP4, progressLabel: "[*] Downloading video as MP4...", fetch: download.DownloadMP4}
	for {
		reader := bufio.NewReader(os.Stdin)
		banner.Clear()
		banner.Show()
		printOptions()

		choice, err := promptLine(reader, "Choose an option: ")
		if err != nil {
			return
		}
		switch choice {
		case "1":
			handleDownload(reader, mp3)
			continue
		case "2":
			handleDownload(reader, mp4)
			continue
		case "3":
			showCredits()
			pressToGoBack()
		default:
			banner.Clear()
			fmt.Println(utils.Gradient("Invalid Option.", utils.Error))
			pressToGoBack()
		}
		banner.Clear()
	}
}

func handleDownload(reader *bufio.Reader, option downloadOption) {
	banner.Clear()
	link, err := promptLine(reader, "Enter YouTube link: ")
	if err != nil {
		return
	}
	fileName, err := promptLine(reader, "Save file as (without extension): ")
	if err != nil {
		return
	}
	downloadsPath := download.DownloadsDir()
	savePath := filepath.Join(downloadsPath, fileName+option.extension)

	fmt.Println(utils.Gradient(option.progressLabel, utils.Minty))

	if err := option.fetch(link, savePath); err != nil {
		reportError(err)
		pressToGoBack()
		return
	}

	reportDownloaded(savePath)
	openDownloadsFolder(downloadsPath)

	fmt.Print(utils.Gradient("\r\nPress ENTER to go back", utils.Minty))
	fmt.Scanln()
	banner.Clear()
}

func printOptions() {
	fmt.Println(utils.Gradient("[1] MP3", utils.Minty))
	fmt.Println(utils.Gradient("[2] MP4", utils.Minty))
	fmt.Println(utils.Gradient("[3] Credits", utils.Minty))
	fmt.Println()
}

func promptLine(reader *bufio.Reader, label string) (string, error) {
	fmt.Print(utils.Gradient(label, utils.Minty) + " ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}
	return strings.TrimSpace(line), nil
}

func reportError(err error) {
	fmt.Println(utils.Gradient(fmt.Sprintf("Error: %v", err), utils.Error))
}

func reportDownloaded(savePath string) {
	fmt.Println(utils.Gradient("\n[✓] Download complete!", utils.Success))
	fmt.Println(utils.Gradient("Saved to: ", utils.Minty) + savePath)
}

func openDownloadsFolder(downloadsPath string) {
	if err := download.OpenDownloadsDir(downloadsPath); err != nil {
		reportError(err)
	}
}

func pressToGoBack() {
	fmt.Println(utils.Gradient("\r\nPress ENTER to go back", utils.Minty))
	fmt.Scanln()
}

func showCredits() {
	banner.Clear()
	banner.Show()
	fmt.Println("Contacts: https://yourpov.dev/")
}
