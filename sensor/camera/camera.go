package camera

import (
	"fmt"
	"os"
	"time"

	"github.com/dhowden/raspicam"
	"github.com/indiependente/motionbot/video/muxer"
	"github.com/pkg/errors"
)

// Camera represents a Camera object which can take pictures and videos.
type Camera interface {
	Picture() (string, error)
	Video() (string, error)
}

// NoIRCamera represents a NoIR Camera.
type NoIRCamera struct {
	Muxer muxer.Muxer
}

// Picture takes a picture and returns the related file and an error.
func (c *NoIRCamera) Picture() (string, error) {
	filename := time.Now().Format("02-01-2006T15:04:05") + ".jpg"
	f, err := os.Create(filename)
	if err != nil {
		return "", errors.Wrapf(err, "Could not create file: %s", filename)
	}
	defer f.Close()

	s := raspicam.NewStill()
	errCh := make(chan error)
	go func() {
		for x := range errCh {
			fmt.Fprintf(os.Stderr, "%v\n", x)
		}
	}()
	raspicam.Capture(s, f, errCh)
	return filename, nil
}

func (c *NoIRCamera) Video() (string, error) {
	filename := time.Now().Format("02012006T150405") + ".h264"
	f, err := os.Create(filename)
	if err != nil {
		return "", errors.Wrapf(err, "Could not create file: %s", filename)
	}
	defer f.Close()

	v := raspicam.NewVid()
	v.Height = 480
	v.Width = 480
	v.Timeout = 10 * time.Second
	errCh := make(chan error)
	go func() {
		for x := range errCh {
			fmt.Fprintf(os.Stderr, "%v\n", x)
		}
	}()
	raspicam.Capture(v, f, errCh)

	mp4Filename, err := c.Muxer.Mux(filename)
	if err != nil {
		return "", errors.Wrap(err, "Could not convert captured video")
	}
	return mp4Filename, nil
}
