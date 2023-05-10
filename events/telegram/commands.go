package telegram

import (
	"errors"
	"log"
	"net/url"
	"note-adviser-bot/events"
	"note-adviser-bot/lib/e"
	"note-adviser-bot/storage"
	"strconv"
	"strings"
)

const NotePrefix = "Note"

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddURLCmd(text) {
		return p.savePage(chatID, text, username)
	}

	if isAddNoteCmd(text) {
		return p.saveNote(chatID, text, username)
	}

	switch text {
	case events.StartCmd:
		// return first layer keyboard
		return p.sendHello(chatID)
	case events.MngLinks, events.MngLinksLabel:
		// retur keyboard
		return p.sendMngLinks(chatID)
	case events.MngNotes, events.MngNotesLabel:
		// return keyboard
		return p.sendMngNotes(chatID)
	case events.HelpSection, events.HelpSectionLabel:
		// return keyboard
		return p.sendHelpSection(chatID)
	case events.ChatInfo, events.ChatInfoLabel:
		// return keyboard
		return p.sendChatInfo(chatID)
	case events.RndLinkCmd, events.RndLinkLabel:
		return p.sendRandomLink(chatID, username)
	case events.ListLinksCmd, events.ListLinksLabel:
		// return list of links
		return p.sendLinksList(chatID)
	case events.DelLinkCmd, events.DelLinkLabel:
		// return delete link
		return p.deleteLink(chatID)
	case events.HelpCmd, events.HelpLabel:
		return p.sendHelp(chatID)
	case events.ContactsCmd, events.ContactsLabel:
		// return contacts
		return p.tg.SendMessage(chatID, events.DevUsername)
	case events.RndNoteCmd, events.RndNoteLabel:
		// return random note
		return p.sendRandomNote(chatID)
	case events.ListNoteCmd, events.ListNoteLabel:
		// return list of notes
		return p.sendNotesList(chatID)
	case events.DelNoteCmd, events.DelNoteLabel:
		return p.deleteNote(chatID)
	case events.GetChatID, events.GetChatIDLabel:
		return p.tg.SendMessage(chatID, strconv.Itoa(chatID))
	case events.GetUsername, events.GetUsernameLabel:
		return p.tg.SendMessage(chatID, username)
	case events.BackCmd, events.BackLabel:
		return p.sendHome(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) saveNote(chatID int, text string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't save note", err) }()

	if err := p.tg.SendMessage(chatID, "Save note: not implemented"); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandomLink(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavePages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send hello message", err) }()

	// if err := p.tg.SendMessage(chatID, msgHello); err != nil {
	// 	return err
	// }

	if err := p.tg.SendKeyboard(chatID, msgHello, events.StartLayer); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendMngLinks(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send manage links keyboard", err) }()

	if err := p.tg.SendKeyboard(chatID, events.MngLinksLabel, events.MngLinksLayer); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendMngNotes(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send manage notes keyboard", err) }()

	if err := p.tg.SendKeyboard(chatID, events.MngNotesLabel, events.MngNotesLayer); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHelpSection(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send help section keyboard", err) }()

	if err := p.tg.SendKeyboard(chatID, events.HelpSectionLabel, events.HelpLayer); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendChatInfo(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send chat info keyboard", err) }()

	if err := p.tg.SendKeyboard(chatID, events.ChatInfoLabel, events.ChatDataLayer); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHome(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send chat info keyboard", err) }()

	if err := p.tg.SendKeyboard(chatID, events.BackLabel, events.StartLayer); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendLinksList(chatID int) error {
	return p.tg.SendMessage(chatID, "Get link list: not implemented")
}

func (p *Processor) deleteLink(chatID int) error {
	return p.tg.SendMessage(chatID, "Delete link: not implemented")
}

func (p *Processor) deleteNote(chatID int) error {
	return p.tg.SendMessage(chatID, "Delete note: not implemented")
}

func (p *Processor) sendRandomNote(chatID int) error {
	return p.tg.SendMessage(chatID, "Get random note: not implemented")
}

func (p *Processor) sendNotesList(chatID int) error {
	return p.tg.SendMessage(chatID, "Get notes list: not implemented")
}

func isAddURLCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func isAddNoteCmd(text string) bool {
	return isNote(text)
}

func isNote(text string) bool {
	return strings.HasPrefix(text, NotePrefix)
}
