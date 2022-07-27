package telegram

import (
	"errors"
	"home/internal/clients/telegram"
	"home/internal/events"
	"home/pkg/lib/e"
	"home/pkg/lib/storage"
)

var (
	ErrUnknownEventType = errors.New("can't process message")
	ErrUnknownMetaType  = errors.New("can't get meta")
)

type Processor struct {
	Tg      *telegram.Client
	Offset  int
	Storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		Tg:      client,
		Storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.Tg.Updates(p.Offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		res = append(res, event(update))
	}

	p.Offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)

	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatId:   u.ID,
			Username: u.Message.From.Username,
		}
	}
	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message != nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message != nil {
		return ""
	}
	return u.Message.Text
}
