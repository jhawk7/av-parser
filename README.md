# AV Parser 
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white) ![YouTube](https://img.shields.io/badge/YouTube-%23FF0000.svg?style=flat&logo=YouTube&logoColor=white) ![FFmpeg](https://shields.io/badge/FFmpeg-%23171717.svg?logo=ffmpeg&style=flat&labelColor=171717&logoColor=5cb85c)

AV Parser is a command-line tool for downloading YouTube videos and extracting audio using [`go-ytdlp](https://github.com/lrstanley/go-ytdlp) and [ffmpeg-go](https://github.com/u2takey/ffmpeg-go). It supports saving both audio and video files, or just one of them, and moves the results to configurable storage directories specified in the environment.

## Features

- Downloads YouTube videos in MP4 format using yt-dlp.
- Extracts audio from downloaded videos as MP3 using ffmpeg.
- Supports flags to save only audio (`-a`) or only video (`-v`).
- Moves processed files to storage directories specified by environment variables.
- Cleans up temporary files after processing.

## Requirements

- Go 1.24.1 or newer
- [`go-ytdlp](https://github.com/lrstanley/go-ytdlp) 
- [ffmpeg-go](https://github.com/u2takey/ffmpeg-go)

## Installation

Clone the repository:

```sh
git clone https://github.com/jhawk7/av-parser.git
cd av-parser
```

Install dependencies:

```sh
go mod tidy
```

## Configuration

Set the following environment variables to specify where audio and video files will be stored. 

```sh
export AV_AUDIO_STORAGE_DIR="/path/to/audio/storage/"
export AV_VIDEO_STORAGE_DIR="/path/to/video/storage/"
```

## Usage

Run the parser with a YouTube URL:

```sh
go run cmd/av-parser/main.go <url> [flag]
```

- `<url>`: The YouTube video URL.
- `[flag]`: Optional. Use `-a` for audio only, `-v` for video only. By default, both are saved.

**Examples:**

Download both audio and video:

```sh
go run cmd/av-parser/main.go https://youtube.com/watch?v=example
```

Download audio only:

```sh
go run cmd/av-parser/main.go https://youtube.com/watch?v=example -a
```

Download video only:

```sh
go run cmd/av-parser/main.go https://youtube.com/watch?v=example -v
```

## Project Structure

```
go.mod
go.sum
README.md
cmd/
  av-parser/
    main.go
```

- [`cmd/av-parser/main.go`](cmd/av-parser/main.go): Main application logic.

## License

MIT License. See [LICENSE](LICENSE) for details.

---
