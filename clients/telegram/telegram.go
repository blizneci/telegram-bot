package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	events "note-adviser-bot/events"
	"note-adviser-bot/lib/e"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

type KeyboardMessage struct {
	ChatID      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// do request <- getUpdates
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendKeyboard(chatID int, text string, layer int) error {
	var msg KeyboardMessage

	msg.ChatID = chatID
	msg.Text = text
	msg.ReplyMarkup = getKeyboard(layer)

	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, sendMessageMethod),
	}
	_, err = http.Post(u.String(), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return e.Wrap("cant't send keyboard", err)
	}

	return nil

}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getKeyboard(layer int) ReplyKeyboardMarkup {
	var k ReplyKeyboardMarkup
	switch layer {
	case events.StartLayer:
		k.Keyboard = [][]KeyboardButton{
			{
				{Text: events.MngLinksLabel},
				{Text: events.MngNotesLabel},
			},
			{
				{Text: events.HelpSectionLabel},
				{Text: events.ChatInfoLabel},
			},
		}
		return k
	case events.MngLinksLayer:
		k.Keyboard = [][]KeyboardButton{
			{
				{Text: events.RndLinkLabel},
				{Text: events.ListLinksLabel},
			},
			{
				{Text: events.DelLinkLabel},
				{Text: events.BackLabel},
			},
		}
	case events.HelpLayer:
		k.Keyboard = [][]KeyboardButton{
			{
				{Text: events.HelpLabel},
				{Text: events.ContactsLabel},
			},
			{
				{Text: events.BackLabel},
			},
		}
	case events.MngNotesLayer:
		k.Keyboard = [][]KeyboardButton{
			{
				{Text: events.RndNoteLabel},
				{Text: events.ListNoteLabel},
			},
			{
				{Text: events.DelNoteLabel},
				{Text: events.BackLabel},
			},
		}
	case events.ChatDataLayer:
		k.Keyboard = [][]KeyboardButton{
			{
				{Text: events.GetChatIDLabel},
				{Text: events.GetUsernameLabel},
			},
			{
				{Text: events.BackLabel},
			},
		}
	}

	return k
}
