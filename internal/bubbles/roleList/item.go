package rolelist

type Item struct {
	title       string
	description string
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	return i.description
}

func (i Item) FilterValue() string {
	return i.title
}
