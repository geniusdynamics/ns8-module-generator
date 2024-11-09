
package commands

import (
	"fmt"
	"log"
	"ns8-module-generator/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type GithubUserName struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputGithubUsername() {
	p := tea.NewProgram(githubUsernameInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(GithubUserName)
	if ok {
		utils.SetGithubUsername(inputModel.value)
	}
}

func githubUsernameInputModel() GithubUserName {
	ti := textinput.New()
	ti.Placeholder = "Github Username"
	ti.Focus()
	ti.CharLimit = 50
	return GithubUserName{
		textInput: ti,
		err:       nil,
	}
}

func (m GithubUserName) Init() tea.Cmd {
	return textinput.Blink
}

func (m GithubUserName) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m GithubUserName) View() string {
	return fmt.Sprintf(
		"\nPersonal Username?\n\n%s",
		m.textInput.View(),
	) + "\n"
}
