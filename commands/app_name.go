package commands

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	textInput textinput.Model
	err       error
}

func InputAppName() {
	p := tea.NewProgram(textInputModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func textInputModel() TextInput {
	ti := textinput.New()
	ti.Placeholder = "App Name"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	return TextInput{
		textInput: ti,
		err:       nil,
	}
}

func (m TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
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

func (m TextInput) View() string {
	return fmt.Sprintf(
		"The name of the app(Should be one word or separated by hyphen)?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
