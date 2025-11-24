package table

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eduardofuncao/pam/internal/db"
)

func Render(tableData *db.TableData, elapsed time.Duration) error {
	model := New(tableData, elapsed)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
