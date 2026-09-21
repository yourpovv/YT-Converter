package download

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	ExtMP3 = ".mp3"
	ExtMP4 = ".mp4"
)

const youtubeExtractorArgs = "youtube:player_client=android,web"

type converterTools struct {
	ytDlp  string
	ffmpeg string
}

func YtDlpPath() (string, error) {
	if path, err := exec.LookPath("yt-dlp"); err == nil {
		return path, nil
	}
	localYtDlp := filepath.Join("utils", "ffmpeg", "yt-dlp.exe")
	if _, err := os.Stat(localYtDlp); err != nil {
		return "", fmt.Errorf("yt-dlp is not installed (no global binary and no bundled copy): install it from https://github.com/yt-dlp/yt-dlp and place the binary in utils/ffmpeg: %w", err)
	}
	return localYtDlp, nil
}

func FfmpegPath() (string, error) {
	absFfmpegPath, err := filepath.Abs(filepath.Join("utils", "ffmpeg", "ffmpeg.exe"))
	if err != nil {
		return "", fmt.Errorf("could not find ffmpeg: %w", err)
	}
	return absFfmpegPath, nil
}

func createDownloadsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("locate working directory: %w", err)
	}
	downloadsPath := filepath.Join(dir, "downloads")
	if err := os.MkdirAll(downloadsPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("create downloads folder %s: %w", downloadsPath, err)
	}
	return downloadsPath, nil
}

func OpenDownloadsDir(path string) error {
	var opener string
	var args []string
	switch runtime.GOOS {
	case "windows":
		opener, args = "explorer", []string{path}
	case "darwin":
		opener, args = "open", []string{path}
	default:
		opener, args = "xdg-open", []string{path}
	}
	if err := exec.Command(opener, args...).Start(); err != nil {
		return fmt.Errorf("open downloads folder %s: %w", path, err)
	}
	return nil
}

func DownloadMP3(youtubeURL, outputPath string) error {
	tools, err := resolveConverterTools()
	if err != nil {
		return err
	}

	absPath, err := resolveSavePath(outputPath)
	if err != nil {
		return err
	}

	if err := runWithOutput(tools.ytDlp, "--extractor-args", youtubeExtractorArgs, "-x", "--audio-format", "mp3", "--ffmpeg-location", tools.ffmpeg, "-o", absPath, youtubeURL); err != nil {
		return fmt.Errorf("download audio: %w", err)
	}

	return tagMP3Title(tools.ffmpeg, absPath, youtubeURL)
}

func DownloadMP4(youtubeURL, outputPath string) error {
	tools, err := resolveConverterTools()
	if err != nil {
		return err
	}

	absPath, err := resolveSavePath(outputPath)
	if err != nil {
		return err
	}

	err = runWithOutput(tools.ytDlp, "--impersonate", "chrome", "-f", "bestvideo[protocol=https]+bestaudio[protocol=https]/best[protocol=https]/best", "--merge-output-format", "mp4", "--ffmpeg-location", tools.ffmpeg, "-o", absPath, youtubeURL)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	return nil
}

func resolveConverterTools() (converterTools, error) {
	ytDlp, err := YtDlpPath()
	if err != nil {
		return converterTools{}, err
	}
	ffmpeg, err := FfmpegPath()
	if err != nil {
		return converterTools{}, err
	}
	return converterTools{ytDlp: ytDlp, ffmpeg: ffmpeg}, nil
}

func tagMP3Title(absFfmpegPath, absPath, youtubeURL string) error {
	title := "Downloaded from: " + youtubeURL + " by YT-Converter"
	if err := runWithOutput(absFfmpegPath, "-i", absPath, "-metadata", "title="+title, "-codec", "copy", absPath+"_temp.mp3"); err != nil {
		return fmt.Errorf("tag mp3 title: %w", err)
	}
	if err := os.Rename(absPath+"_temp.mp3", absPath); err != nil {
		return fmt.Errorf("replace tagged mp3: %w", err)
	}
	return nil
}

func resolveSavePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("could not find downloads folder: %w", err)
	}
	return absPath, nil
}

func runWithOutput(binaryPath string, args ...string) error {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
