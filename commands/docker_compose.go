package commands

import (
	"errors"
	"fmt"
	"ns8-module-generator/utils"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FilePicker struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m FilePicker) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
		return m, tea.Quit
	}
	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilePicker) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	// s.WriteString("\n\n press (esc) or backspace to go back")
	return s.String()
}

func PickFile() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".yaml"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := FilePicker{
		filepicker: fp,
	}
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(FilePicker)
	fmt.Print("\033[H\033[2J")
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
	utils.SetDockerComposePath(mm.selectedFile)
}
