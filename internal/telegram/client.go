package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/geqo/togram/internal/detect"
)

const apiBase = "https://api.telegram.org/bot"

type Client struct {
	token string
	http  *http.Client
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{}}
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

func (c *Client) Send(chatID string, t detect.Type, r io.Reader, filename string) error {
	if t == detect.TypeText {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		text := strings.TrimRight(string(data), "\n")
		if len(text) <= 4096 {
			return c.sendMessage(chatID, text)
		}
		// text too long — send as document
		r = strings.NewReader(text)
		if filename == "" {
			filename = "message.txt"
		}
		t = detect.TypeDocument
	}

	if filename == "" {
		filename = "file"
	}
	return c.sendFile(chatID, t, r, filename)
}

func (c *Client) sendMessage(chatID, text string) error {
	params := url.Values{}
	params.Set("chat_id", chatID)
	params.Set("text", text)

	resp, err := c.http.PostForm(c.url("sendMessage"), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkResponse(resp.Body)
}

func (c *Client) sendFile(chatID string, t detect.Type, r io.Reader, filename string) error {
	method, field := methodAndField(t, filename)

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		fw, err := mw.CreateFormFile(field, filepath.Base(filename))
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			pw.CloseWithError(err)
			return
		}
		mw.WriteField("chat_id", chatID)
		pw.CloseWithError(mw.Close())
	}()

	resp, err := c.http.Post(c.url(method), mw.FormDataContentType(), pr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkResponse(resp.Body)
}

func (c *Client) url(method string) string {
	return apiBase + c.token + "/" + method
}

func methodAndField(t detect.Type, filename string) (method, field string) {
	switch t {
	case detect.TypePhoto:
		return "sendPhoto", "photo"
	case detect.TypeVideo:
		return "sendVideo", "video"
	case detect.TypeAudio:
		return "sendAudio", "audio"
	default:
		return "sendDocument", "document"
	}
}

func checkResponse(body io.Reader) error {
	var r apiResponse
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return err
	}
	if !r.OK {
		return fmt.Errorf("telegram: %s", r.Description)
	}
	return nil
}
