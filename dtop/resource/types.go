package resource

import "github.com/charmbracelet/bubbles/list"

type ResourceItem struct {
	Id       string
	Title    string
	Desc     string
	Selected bool
}

func (i ResourceItem) TitleFunc() string {
	prefix := "[ ]"
	if i.Selected {
		prefix = "[x]" // checkmark
	}
	return prefix + " " + i.Title
}

func (i ResourceItem) Description() string { return i.Desc }
func (i ResourceItem) FilterValue() string { return i.Title }
func (i ResourceItem) ID() string          { return i.Id }

// Messages pour les opérations asynchrones
type ResourcesLoadedMsg struct {
	Items []list.Item
	Err   error
}

type ResourcesDeletedMsg struct{ Errs []error }

type ImagesPublishedMsg struct {
	Errs []error
}