package convert

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func Convert(path string, w io.Writer) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", path,
		"-vn",                 // drop any video streams (covers)
		"-map_metadata", "-1", // strip all metadata
		"-f", "mp3", // convert to mp3
		"-", // output to stdout
	)

	var stderr bytes.Buffer

	cmd.Stdout = w
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec ffmpeg: %v: %s", err, stderr.Bytes())
	}

	return nil
}
