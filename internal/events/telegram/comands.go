package telegram

import (
	"errors"
	"home/pkg/lib/e"
	"home/pkg/lib/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command \"%s\" from [%s]", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// start: /start: helloMess + helpInfo

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.Tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExist, err := p.Storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.Tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.Storage.Save(page); err != nil {
		return err
	}

	if err := p.Tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't random", err) }()

	page, err := p.Storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.Tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.Tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}

	return p.Storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.Tg.SendMessage(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.Tg.SendMessage(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
