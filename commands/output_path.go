package commands

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type OutputPathInputText struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputOutputDirPath() (string, error) {
	p := tea.NewProgram(outputPathInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(OutputPathInputText)
	if ok {
		return inputModel.value, nil
	}
	return "", fmt.Errorf("Error occurred while reading output dir")
}

func outputPathInputModel() OutputPathInputText {
	ti := textinput.New()
	ti.Placeholder = "generated-module"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.SetValue("generated-module") // Set default value
	return OutputPathInputText{
		textInput: ti,
		err:       nil,
	}
}

func (m OutputPathInputText) Init() tea.Cmd {
	return textinput.Blink
}

func (m OutputPathInputText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			currentValue := m.textInput.Value()
			if currentValue == "" {
				m.err = fmt.Errorf("Output path cannot be empty.")
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

func (m OutputPathInputText) View() string {
	return fmt.Sprintf(
		"\nEnter the relative output path(use `pwd` )?\n\n%s",
		m.textInput.View(),
	) + "\n"
}
