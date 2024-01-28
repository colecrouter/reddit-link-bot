package video

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

const testAudio = "https://v.redd.it/m7gfgtu9sebb1/DASH_audio.mp4"
const testVideo = "https://v.redd.it/m7gfgtu9sebb1/DASH_360.mp4"

func TestMerge(t *testing.T) {
	urlAudio, _ := url.Parse(testAudio)
	urlVideo, _ := url.Parse(testVideo)

	audioResp, _ := http.Get(urlAudio.String())
	videoResp, _ := http.Get(urlVideo.String())

	v, err := Merge(audioResp.Body, videoResp.Body)
	if err != nil {
		t.Error(err)
	} else if v == nil {
		t.Error("video is nil")
	}

	bytes, err := io.ReadAll(v)
	if err != nil {
		t.Error(err)
	}

	if len(bytes) == 0 {
		t.Error("video is empty")
	}
}
