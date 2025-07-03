package commands

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type GithubTokenInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputGithubToken() (string, error) {
	p := tea.NewProgram(githubTokenInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(GithubTokenInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("error reading the github token")
}

func githubTokenInputModel() GithubTokenInput {
	ti := textinput.New()
	ti.Placeholder = "Github Token"
	ti.Focus()
	ti.CharLimit = 255
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	ti.SetValue("") // Set default value to empty
	return GithubTokenInput{
		textInput: ti,
		err:       nil,
	}
}

func (m GithubTokenInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m GithubTokenInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
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

func (m GithubTokenInput) View() string {
	return fmt.Sprintf(
		"\nGithub token used for pushing the module to github?\n\n%s",
		m.textInput.View(),
	) + "\n"
}
