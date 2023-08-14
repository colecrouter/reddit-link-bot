package video

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
)

func Merge(audio string, video string) (r *io.ReadCloser, err error) {
	// We can actually pass URLs directly to ffmpeg, but that requires a special
	// build of ffmpeg with HTTPS enabled. Instead, we'll download the files manually

	var err1, err2 error
	var videoData, audioData *io.ReadCloser
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		videoData, err1 = fetch.Fetch(video)
		wg.Done()
	}()

	go func() {
		audioData, err2 = fetch.Fetch(audio)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		err = err1
		return
	} else if err2 != nil {
		err = err2
		return
	}

	// Create temp files to store the downloaded data. Named pipes don't exist on Windows
	vFile, _ := os.CreateTemp("", "video_*")
	aFile, _ := os.CreateTemp("", "audio_*")
	defer vFile.Close()
	defer aFile.Close()

	io.Copy(vFile, *videoData)
	io.Copy(aFile, *audioData)

	cmd := exec.Command("ffmpeg", "-y", "-i", vFile.Name(), "-i", aFile.Name(), "-map", "0:0", "-map", "1:0", "-f", "ismv", "-c:v", "copy", "pipe:")
	// cmd.Stderr = os.Stderr
	output, _ := cmd.StdoutPipe()

	err = cmd.Start()
	if err != nil {
		return
	}

	// After the command has finished, delete the old input files.
	// There might be a better way to handle this, such that the files get deleted if cmd.Start() fails
	go func() {
		cmd.Wait()
		os.Remove(vFile.Name())
		os.Remove(aFile.Name())
	}()

	// Read the file into memory
	r = &output

	return
}
