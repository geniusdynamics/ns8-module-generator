package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputAppName() (string, error) {
	p := tea.NewProgram(textInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(TextInput)
	if !ok {
		return "", fmt.Errorf("Could not cast bubble teamodel")
	}
	return inputModel.value, nil
}

func textInputModel() TextInput {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-app"
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
		case tea.KeyEnter:
			currentValue := m.textInput.Value()
			if strings.Contains(currentValue, " ") {
				m.err = fmt.Errorf("App name cannot contain spaces. Use hyphens instead.")
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

func (m TextInput) View() string {
	return fmt.Sprintf(
		"The name of the app(Should be one word or separated by hyphen)?\n\n%s",
		m.textInput.View(),
	) + "\n"
}
