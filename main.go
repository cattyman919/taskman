package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/windows"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type model struct {
	processes table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.processes.Focused() {
				m.processes.Blur()
			} else {
				m.processes.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.processes.SelectedRow()[1]),
			)
		}
	}
	m.processes, cmd = m.processes.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.processes.View()) + "\n"
}

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	columns := []table.Column{
		{Title: "PID", Width: 10},
		{Title: "Name", Width: 20},
	}

	process_list, _ := process.Processes()

	name_process := []table.Row{}

	for _, p := range process_list {
		if name, err := p.Name(); err == nil {
			pid, _ := p.Ppid()

			if p.Pid == 0 || name == "System" {
				continue
			}

			name_process = append(name_process, table.Row{strconv.Itoa(int(pid)), name})

		} else {
			log.Print(err)
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(name_process),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	p := tea.NewProgram(model{
		processes: t,
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func amAdmin() bool {
	elevated := windows.GetCurrentProcessToken().IsElevated()
	return elevated
}
