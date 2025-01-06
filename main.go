package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/windows"
)

const PROCESS_LIMIT int = 10

type model struct {
	processes []string
	cursor    int
	selected  map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < PROCESS_LIMIT-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {

	s := "Processes: \n\n"
	shown, i := 0, 0

	for shown < PROCESS_LIMIT {
		name := m.processes[i]
		if name == "" {
      i++
			continue
		}

		cursor := " " // no cursor
		if m.cursor == shown {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if _, ok := m.selected[shown]; ok {
			checked = "x" // selected!
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, name)

    i++
		shown++
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	process, _ := process.Processes()

	name_process := make([]string, len(process))

	for i, p := range process {
		if name, err := p.Name(); err == nil {
			name_process[i] = name
		} else {
			log.Print(err)
		}
	}

	process = nil

	p := tea.NewProgram(model{
		processes: name_process,
		selected:  make(map[int]struct{}),
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
