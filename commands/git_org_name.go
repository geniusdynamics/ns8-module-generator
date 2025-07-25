package commands

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type GitOrgNameTextInput struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputGithubOrganizationName() (string, error) {
	p := tea.NewProgram(githubTextInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(GitOrgNameTextInput)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("an error occurred while reading github organization")
}

func githubTextInputModel() GitOrgNameTextInput {
	ti := textinput.New()
	ti.Placeholder = "Organization Github Username"
	ti.Focus()
	ti.CharLimit = 50
	ti.SetValue("") // Set default value to empty
	return GitOrgNameTextInput{
		textInput: ti,
		err:       nil,
	}
}

func (m GitOrgNameTextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m GitOrgNameTextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m GitOrgNameTextInput) View() string {
	return fmt.Sprintf(
		"\nGithub Organization Username?(Leave blank if you need to push to personal account)\n\n%s",
		m.textInput.View(),
	) + "\n"
}
