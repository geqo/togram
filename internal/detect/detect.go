package detect

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type Type int

const (
	TypeText Type = iota
	TypePhoto
	TypeVideo
	TypeAudio
	TypeDocument
)

func FromReader(r io.Reader) (Type, io.Reader, error) {
	buf := make([]byte, 512)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return TypeDocument, nil, err
	}
	buf = buf[:n]

	mime := http.DetectContentType(buf)
	rest := io.MultiReader(bytes.NewReader(buf), r)
	return fromMIME(mime), rest, nil
}

func FromFilename(name string) Type {
	lower := strings.ToLower(name)
	switch {
	case hasSuffix(lower, ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"):
		return TypePhoto
	case hasSuffix(lower, ".mp4", ".avi", ".mov", ".mkv", ".webm"):
		return TypeVideo
	case hasSuffix(lower, ".mp3", ".ogg", ".flac", ".wav", ".m4a"):
		return TypeAudio
	default:
		return TypeDocument
	}
}

func fromMIME(mime string) Type {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return TypePhoto
	case strings.HasPrefix(mime, "video/"):
		return TypeVideo
	case strings.HasPrefix(mime, "audio/"):
		return TypeAudio
	case strings.HasPrefix(mime, "text/"):
		return TypeText
	default:
		return TypeDocument
	}
}

func hasSuffix(s string, suffixes ...string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}
