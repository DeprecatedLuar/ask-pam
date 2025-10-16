package table

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func Render(columns []string, data [][]string, elapsed time.Duration) error {
	model := New(columns, data, elapsed)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
