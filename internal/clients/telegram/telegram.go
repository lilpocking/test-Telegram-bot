package telegram

import (
	"encoding/json"
	"errors"
	"home/pkg/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	Host     string
	BasePath string
	Client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		Host:     host,
		BasePath: newBasePath(token),
		Client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// Get не пишется
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.WrapIfErr("can't convert resoonse json data in structer: %w", err)
	}
	if !res.OK {
		return nil, errors.New("response message is not ok")
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.WrapIfErr("can't send message: %w", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	reqErr := "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, e.WrapIfErr(reqErr, err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.Client.Do(req)

	if err != nil {
		return nil, e.WrapIfErr(reqErr, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.WrapIfErr(reqErr, err)
	}

	return body, nil
}
