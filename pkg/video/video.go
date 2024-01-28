package video

import (
	"io"
	"os"
	"os/exec"
)

func Merge(audio io.ReadCloser, video io.ReadCloser) (r io.ReadCloser, err error) {
	// Create temp files to store the downloaded data. Named pipes don't exist on Windows
	vFile, _ := os.CreateTemp("", "video_*")
	aFile, _ := os.CreateTemp("", "audio_*")
	defer vFile.Close()
	defer aFile.Close()

	io.Copy(vFile, video)
	io.Copy(aFile, audio)

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
	r = output

	return
}
