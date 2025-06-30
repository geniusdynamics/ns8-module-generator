package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AppInitGitTextInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputAppGitInit() (string, error) {
	p := tea.NewProgram(appGitInitModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(AppInitGitTextInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("an error occurred while initialising git")
}

func appGitInitModel() AppInitGitTextInput {
	ti := textinput.New()
	ti.Placeholder = "yes / no"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.SetValue("no") // Set default value
	return AppInitGitTextInput{
		textInput: ti,
		err:       nil,
	}
}

func (m AppInitGitTextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m AppInitGitTextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			currentValue := strings.ToLower(m.textInput.Value())
			if currentValue != "yes" && currentValue != "no" {
				m.err = fmt.Errorf("Please enter 'yes' or 'no'.")
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

func (m AppInitGitTextInput) View() string {
	return fmt.Sprintf(
		"\nInitialise and Push to Git(yes or no)\n\n%s",
		m.textInput.View(),
	) + "\n"
}
