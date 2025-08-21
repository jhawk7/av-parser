package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lrstanley/go-ytdlp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	TMP_VID_FOLDER   = "./yt-tmp/"
	TMP_AUDIO_FOLDER = "./audio-tmp/"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing url; Usage: 'go run main.go <url> <optional '-s' flag to save video>'")
		return
	}

	ytUrl := os.Args[1]
	saveFlag := false
	if len(os.Args) == 3 && strings.ToLower(os.Args[2]) == "-s" {
		saveFlag = true
	}

	downloadContent(ytUrl)
	parseAV()
	transferFiles(saveFlag)
}

func downloadContent(ytUrl string) {
	fmt.Printf("retrieving video file from url %s\n", ytUrl)

	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)

	if osErr := os.MkdirAll(TMP_VID_FOLDER, os.ModePerm); osErr != nil {
		err := fmt.Errorf("failed to make tmp video dir; [err: %v]", osErr)
		errorHandler(err, true)
		return
	}

	dl := ytdlp.New().
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output(TMP_VID_FOLDER + "%(extractor)s - %(title)s.%(ext)s")

	_, dlErr := dl.Run(context.TODO(), ytUrl)
	if dlErr != nil {
		err := fmt.Errorf("failed to download content; [err: %v]", dlErr)
		errorHandler(err, true)
	}

	fmt.Println("vid download complete.")
}

func parseAV() {
	fmt.Println("parsing audio from video..")

	if osErr := os.MkdirAll(TMP_AUDIO_FOLDER, os.ModePerm); osErr != nil {
		err := fmt.Errorf("failed to make tmp audio dir; [err: %v]", osErr)
		errorHandler(err, true)
		return
	}

	entries, osErr := os.ReadDir(TMP_VID_FOLDER)
	if osErr != nil {
		err := fmt.Errorf("failed to read dir; [err: %v]", osErr)
		errorHandler(err, true)
		return
	}

	if len(entries) == 0 {
		err := fmt.Errorf("no files in dir")
		errorHandler(err, true)
		return
	}

	filename := entries[0].Name()
	audioFilename := strings.Split(filename, ".")[0] + ".mp3"

	streamErr := ffmpeg.Input(fmt.Sprintf("%s%s", TMP_VID_FOLDER, filename)).
		Output(fmt.Sprintf("%s%s", TMP_AUDIO_FOLDER, audioFilename), ffmpeg.KwArgs{"q:a": 0, "map": "a"}).
		OverWriteOutput().ErrorToStdOut().Run()

	if streamErr != nil {
		err := fmt.Errorf("failed to parse video to audio file; [err: %v]", streamErr)
		errorHandler(err, true)
		return
	}

	fmt.Println("av parsing complete.")
}

func transferFiles(flag bool) {
	vidArchive, isSet := os.LookupEnv("AV_VIDEO_STORAGE_DIR")
	if !isSet {
		err := fmt.Errorf("video storage dir not set")
		errorHandler(err, true)
	}

	audoArchive, isSet := os.LookupEnv("AV_AUDIO_STORAGE_DIR")
	if !isSet {
		err := fmt.Errorf("audio storage dir not set")
		errorHandler(err, true)
	}

	copyFiles := func(srcFile, dstDir string) {
		entries, _ := os.ReadDir(srcFile)
		for _, entry := range entries {
			src, srcErr := os.Open(fmt.Sprintf("%s%s", srcFile, entry.Name()))
			errorHandler(srcErr, true)
			dst, dstErr := os.Create(dstDir + entry.Name())
			errorHandler(dstErr, true)

			defer func() {
				src.Close()
				dst.Close()
			}()

			if _, cpErr := io.Copy(dst, src); cpErr != nil {
				err := fmt.Errorf("failed to copy files; [err: %v]", cpErr)
				errorHandler(err, true)
			}
		}
	}

	fmt.Println("transferring audio files..")
	copyFiles(TMP_AUDIO_FOLDER, audoArchive)
	if flag {
		fmt.Println("transferring video files..")
		copyFiles(TMP_VID_FOLDER, vidArchive)
	}
	fmt.Println("transfer complete")

	//removing tmp folders
	os.RemoveAll(TMP_VID_FOLDER)
	os.RemoveAll(TMP_AUDIO_FOLDER)
}

func errorHandler(err error, fatal bool) {
	if err != nil {
		if fatal {
			panic(err)
		} else {
			fmt.Printf("error %v", err)
		}
	}
}
