package telegram

const msgHelp = `I can save and keep your pages and notes.

Also I can offer you them to read.

In order to save the page, just send me a link to it.
In order to save the note, just add "Note" before the text and send it to me.

In order to get a random page from your list, use keyboard or send me command /rndlink.
In order to get a random note from your list, use keyboard or send me command /rndnote.
Caution! After that, this page/note will be removed from your list!`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages/notes 🙊"
	msgSaved          = "Saved 👌"
	msgAlreadyExists  = "You have already have this page/note in your list 🤗"
)
