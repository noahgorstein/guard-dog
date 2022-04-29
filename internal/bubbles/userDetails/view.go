package userdetails

func (b Bubble) View() string {

	if b.active {
		b.viewport.Style = b.Styles.ActiveViewportStyle
	}
	return b.viewport.View()
}
