package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

const DevUsername = "https://t.me/Zhelvakov"

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

const (
	StartLayer int = iota
	MngLinksLayer
	MngNotesLayer
	HelpLayer
	ChatDataLayer
)

const (
	MngLinksLabel    = "Управление ссылками"
	MngNotesLabel    = "Управление заметками"
	HelpSectionLabel = "Помощь"
	ChatInfoLabel    = "Данные чата"
	RndLinkLabel     = "Получить случайную ссылку"
	ListLinksLabel   = "Получить список ссылок"
	DelLinkLabel     = "Удалить ссылку"
	HelpLabel        = "Справка"
	ContactsLabel    = "Контакты"
	RndNoteLabel     = "Получить случайную заметку"
	ListNoteLabel    = "Получить список заметок"
	DelNoteLabel     = "Удалить заметку"
	GetChatIDLabel   = "Получить ID чата"
	GetUsernameLabel = "Получить Username"
	BackLabel        = "Назад"
)

const (
	StartCmd     = "/start"
	MngLinks     = "/managelinks"
	MngNotes     = "/managenotes"
	HelpSection  = "/helpsection"
	ChatInfo     = "/chatinfo"
	RndLinkCmd   = "/rndlink"
	ListLinksCmd = "/listlinks"
	DelLinkCmd   = "/deletelink"
	HelpCmd      = "/help"
	ContactsCmd  = "/contacts"
	RndNoteCmd   = "/rndnote"
	ListNoteCmd  = "/listnotes"
	DelNoteCmd   = "/deletenote"
	GetChatID    = "/getchatid"
	GetUsername  = "/getusername"
	BackCmd      = "/back"
)
