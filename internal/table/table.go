package table

import tea "github.com/charmbracelet/bubbletea"

func Render(columns []string, data [][]string) error {
	model := New(columns, data)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
