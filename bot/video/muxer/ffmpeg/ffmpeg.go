package ffmpeg

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	ffmpeg         = "ffmpeg"
	framerateParam = "-framerate"
	inputParam     = "-i"
	copyParam      = "-c"
	copyTypeParam  = "copy"
)

type Muxer struct {
	Rate int
}

func (m Muxer) Mux(H264File string) (string, error) {
	filenameMP4 := m.getFileNameMP4(H264File)
	cmd := exec.Command(ffmpeg, m.buildParams(H264File, filenameMP4)...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "Could not re-mux H264 raw file %s into MP4 container %s: %s", H264File, filenameMP4, string(stdoutStderr))
	}
	return filenameMP4, nil
}

func (m Muxer) buildParams(H264File, fnameMP4 string) []string {
	framerate := strconv.Itoa(m.Rate)
	return []string{framerateParam, framerate, inputParam, H264File, copyParam, copyTypeParam, fnameMP4}
}

func (m Muxer) getFileNameMP4(H264File string) string {
	basename := filepath.Base(H264File)
	return strings.TrimSuffix(basename, filepath.Ext(basename)) + ".mp4"
}
