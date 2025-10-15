package table

import (
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) moveUp() Model {
	if m.selectedRow > 0 {
		m.selectedRow--
		if m.selectedRow < m.offsetY {
			m.offsetY = m.selectedRow
		}
	}
	return m
}

func (m Model) moveDown() Model {
	if m.selectedRow < m.numRows()-1 {
		m.selectedRow++
		if m.selectedRow >= m.offsetY+m.visibleRows {
			m.offsetY = m.selectedRow - m.visibleRows + 1
		}
	}
	return m
}

func (m Model) moveLeft() Model {
	if m.selectedCol > 0 {
		m.selectedCol--
		if m.selectedCol < m.offsetX {
			m.offsetX = m.selectedCol
		}
	}
	return m
}

func (m Model) moveRight() Model {
	if m.selectedCol < m.numCols()-1 {
		m.selectedCol++
		if m.selectedCol >= m.offsetX+m.visibleCols {
			m.offsetX = m.selectedCol - m.visibleCols + 1
		}
	}
	return m
}

func (m Model) jumpToFirstCol() Model {
	m.selectedCol = 0
	m.offsetX = 0
	return m
}

func (m Model) jumpToLastCol() Model {
	m.selectedCol = m.numCols() - 1
	if m.visibleCols < m.numCols() {
		m.offsetX = m.numCols() - m.visibleCols
	}
	return m
}

func (m Model) jumpToFirstRow() Model {
	m.selectedRow = 0
	m.offsetY = 0
	return m
}

func (m Model) jumpToLastRow() Model {
	m.selectedRow = m.numRows() - 1
	m.offsetY = m.numRows() - m.visibleRows
	return m
}

func (m Model) pageUp() Model {
	m.selectedRow -= m.visibleRows
	if m.selectedRow < 0 {
		m.selectedRow = 0
	}
	m.offsetY = m.selectedRow
	return m
}

func (m Model) pageDown() Model {
	m.selectedRow += m.visibleRows
	if m.selectedRow >= m.numRows() {
		m.selectedRow = m.numRows() - 1
	}
	if m.selectedRow >= m.offsetY+m.visibleRows {
		m.offsetY = m.selectedRow - m.visibleRows + 1
	}
	return m
}

func (m Model) copySelectedCell() (Model, tea.Cmd) {
	if m.selectedRow >= 0 && m.selectedRow < m.numRows() &&
		m.selectedCol >= 0 && m.selectedCol < m.numCols() {
		go clipboard.WriteAll(m.data[m.selectedRow][m.selectedCol])
		m.blinkCopiedCell = true
		return m, tea.Tick(time.Millisecond*400, func(time.Time) tea.Msg {
			return blinkMsg{}
		})
	}
	return m, nil
}
