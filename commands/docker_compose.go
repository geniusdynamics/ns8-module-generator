package commands

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
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

func PickFile() (string, error) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".yaml", ".yml"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := FilePicker{
		filepicker: fp,
	}
	tm, err := tea.NewProgram(&m).Run()
	mm := tm.(FilePicker)
	fmt.Print("\033[H\033[2J")
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
	if err != nil {
		return "", fmt.Errorf("An error occurred while selecting file path")
	}
	return mm.selectedFile, nil
}

type DockerComposeSourceInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputDockerComposeSource() (string, error) {
	p := tea.NewProgram(dockerComposeSourceInputModel())
	input, err := p.Run()
	if err != nil {
		return "", err
	}
	inputModel, ok := input.(DockerComposeSourceInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("error reading docker compose source")
}

func dockerComposeSourceInputModel() DockerComposeSourceInput {
	ti := textinput.New()
	ti.Placeholder = "local / remote"
	ti.Focus()
	ti.CharLimit = 10
	ti.SetValue("local") // Set default value
	return DockerComposeSourceInput{
		textInput: ti,
		err:       nil,
	}
}

func (m DockerComposeSourceInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m DockerComposeSourceInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			currentValue := strings.ToLower(m.textInput.Value())
			if currentValue != "local" && currentValue != "remote" {
				m.err = fmt.Errorf("Please enter 'local' or 'remote'.")
				return m, nil
			}
			m.value = currentValue
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.value = m.textInput.Value()
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m DockerComposeSourceInput) View() string {
	return fmt.Sprintf(
		"\nIs the Docker Compose file local or remote? (local/remote)\n\n%s",
		m.textInput.View(),
	) + "\n"
}

type DockerComposeUrlInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputDockerComposeUrl() (string, error) {
	p := tea.NewProgram(dockerComposeUrlInputModel())
	input, err := p.Run()
	if err != nil {
		return "", err
	}
	inputModel, ok := input.(DockerComposeUrlInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("error reading docker compose URL")
}

func dockerComposeUrlInputModel() DockerComposeUrlInput {
	ti := textinput.New()
	ti.Placeholder = "Docker Compose URL"
	ti.Focus()
	ti.CharLimit = 255
	return DockerComposeUrlInput{
		textInput: ti,
		err:       nil,
	}
}

func (m DockerComposeUrlInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m DockerComposeUrlInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			currentValue := m.textInput.Value()
			_, err := url.ParseRequestURI(currentValue)
			if err != nil {
				m.err = fmt.Errorf("Invalid URL: %v", err)
				return m, nil
			}
			m.value = currentValue
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.value = m.textInput.Value()
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m DockerComposeUrlInput) View() string {
	return fmt.Sprintf(
		"\nEnter the URL of the Docker Compose file:\n\n%s",
		m.textInput.View(),
	) + "\n"
}
