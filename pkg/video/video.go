package video

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Merge(audio io.ReadCloser, video io.ReadCloser) (r io.ReadCloser, err error) {
	// Create temp files to store the downloaded data. Named pipes don't exist on Windows
	vFile, err := os.CreateTemp("", "video_*")
	if err != nil {
		return nil, fmt.Errorf("create temp video file: %w", err)
	}
	aFile, err := os.CreateTemp("", "audio_*")
	if err != nil {
		vFile.Close()
		os.Remove(vFile.Name())
		return nil, fmt.Errorf("create temp audio file: %w", err)
	}
	// Ensure temporary files are closed; they will be removed after ffmpeg completes
	defer vFile.Close()
	defer aFile.Close()

	if _, err := io.Copy(vFile, video); err != nil {
		return nil, fmt.Errorf("write video temp: %w", err)
	}
	if _, err := io.Copy(aFile, audio); err != nil {
		return nil, fmt.Errorf("write audio temp: %w", err)
	}
	// We can close the input streams as we've consumed them fully.
	video.Close()
	audio.Close()

	cmd := exec.Command("ffmpeg", "-y", "-i", vFile.Name(), "-i", aFile.Name(), "-map", "0:0", "-map", "1:0", "-f", "ismv", "-c:v", "copy", "pipe:")
	cmd.Stderr = os.Stderr
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("ffmpeg start: %w", err)
	}

	// After the command has finished, delete the old input files.
	// There might be a better way to handle this, such that the files get deleted if cmd.Start() fails
	go func() {
		_ = cmd.Wait()
		_ = os.Remove(vFile.Name())
		_ = os.Remove(aFile.Name())
	}()

	// Return a read-closer tied to ffmpeg's stdout
	r = output
	return
}
