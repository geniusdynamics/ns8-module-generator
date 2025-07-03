package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type GitAuthMethodInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputGitAuthMethod() (string, error) {
	p := tea.NewProgram(gitAuthMethodModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(GitAuthMethodInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("error reading the git auth method")
}

func gitAuthMethodModel() GitAuthMethodInput {
	ti := textinput.New()
	ti.Placeholder = "ssh / token"
	ti.Focus()
	ti.CharLimit = 10
	ti.SetValue("token") // Set default value
	return GitAuthMethodInput{
		textInput: ti,
		err:       nil,
	}
}

func (m GitAuthMethodInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m GitAuthMethodInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			currentValue := strings.ToLower(m.textInput.Value())
			if currentValue != "ssh" && currentValue != "token" {
				m.err = fmt.Errorf("Please enter 'ssh' or 'token'.")
				return m, nil
			}
			m.value = currentValue
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.value = m.textInput.Value()
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m GitAuthMethodInput) View() string {
	return fmt.Sprintf("Use SSH or Token for git authentication? %s", m.textInput.View())
}
