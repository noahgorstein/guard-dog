package selector

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name string
}

type Model struct {
	active        bool
	choiceMap     map[string]interface{}
	choices       []string
	cursor        int
	selected      map[int]item
	selectedIndex int
	prompt        string
	height        int
	width         int
	Styles        Styles
}

func New(prompt string, choiceMap map[string]interface{}, height int, width int) Model {
	choices := []string{}
	for k := range choiceMap {
		choices = append(choices, k)
	}

	return Model{
		Styles:    DefaultStyles(),
		prompt:    prompt,
		choiceMap: choiceMap,
		choices:   choices,
		selected:  make(map[int]item),
		width:     width,
		height:    height,
	}
}

func (m Model) GetSelected() interface{} {
	key := m.selected[m.selectedIndex].name
	return m.choiceMap[key]
}

func (m *Model) GetChoices() []string {
	return m.choices
}

func (m *Model) SetChoices(choiceMap map[string]interface{}) {
	if choiceMap == nil {
		choiceMap = make(map[string]interface{})
	}
	choices := []string{}
	for k := range choiceMap {
		choices = append(choices, k)
	}
	m.choiceMap = choiceMap
	m.choices = choices
}

func (m *Model) Reset() {
	m.cursor = 0
	m.selectedIndex = 0
	m.selected = make(map[int]item)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if m.active {

			switch msg.String() {

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			case " ", "enter":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected = map[int]item{}
					m.selectedIndex = m.cursor
					m.selected[m.cursor] = item{name: m.choices[m.cursor]}
				}
			}
		}
	}
	return m, nil
}

func (m *Model) SetIsActive(active bool) {
	m.active = active
}

func (m *Model) Height() int {
	return m.height
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) SetHeight(height int) {
	m.setSize(m.width, height)
}

func (m *Model) SetWidth(width int) {
	m.setSize(width, m.height)
}

func (m *Model) SetSize(width, height int) {
	m.setSize(width, height)
}

func (m *Model) setSize(width, height int) {
	m.width = width
	m.height = height
	m.Styles.ActiveStyle.Width(width).Height(height)
	m.Styles.InactiveStyle.Width(width).Height(height)
}

func (m Model) View() string {
	s := ""

	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = m.Styles.SelectedChoiceStyle.Render(choice)
		}

		// Is this choice selected?
		checked := " "
		if _, ok := m.selected[i]; ok {
			choice = m.Styles.CheckedChoiceStyle.Render(choice)
			checked = m.Styles.CheckmarkStyle.Render("✓")
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
	}

	if m.active {
		return m.Styles.ActiveStyle.Render(m.Styles.PromptStyle.Render(m.prompt) + "\n\n" + s)
	}
	return m.Styles.InactiveStyle.Render(m.Styles.PromptStyle.Render(m.prompt) + "\n\n" + s)
}
