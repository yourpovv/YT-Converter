<div align="center">
<img width="86" src="assets/icon.png" alt="Icon">

# YT-Converter

**Terminal YouTube downloader for MP3 audio and MP4 video**

<img width="995" height="538" alt="image" src="https://github.com/user-attachments/assets/54451a89-de9e-4463-957e-0e655e80f389" />

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows-lightgrey)](https://github.com/yourpovv/YT-Converter)
[![Language](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![GitHub](https://img.shields.io/github/stars/yourpovv/YT-Converter?style=social)](https://github.com/yourpovv/YT-Converter)

</div>

## Features

- Download any YouTube link and convert it to MP3
- Download full YouTube videos as MP4 at the highest available quality
- Tags each MP3 with its source link
- Opens the downloads folder when finished

## Requirements

- [Go 1.23+](https://go.dev/dl)
- `yt-dlp` (global install, or `utils/ffmpeg/yt-dlp.exe`)
- `ffmpeg` at `utils/ffmpeg/ffmpeg.exe`

## Installation

```bash
git clone https://github.com/yourpovv/YT-Converter.git
cd YT-Converter
```

## Usage

```bash
go run .
```

Or build it first:

```bash
go build -o YT-Converter .
./YT-Converter
```

Menu options:

| Option | Action |
| ------ | ------ |
| `1` | Download MP3, enter a YouTube link and a file name |
| `2` | Download MP4, enter a YouTube link and a file name |
| `3` | Show credits |

Files are saved to `downloads/` as `<name>.mp3` or `<name>.mp4`.

## License

[MIT](LICENSE) © [YourPOV](https://github.com/yourpovv)
