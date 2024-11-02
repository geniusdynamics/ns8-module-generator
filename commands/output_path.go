package commands

import (
	"fmt"
	"log"
	"ns8-module-generator/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type OutputPathInputText struct {
	textInput textinput.Model
	err       error
	value     string
}

func InputOutputDirPath() {
	p := tea.NewProgram(outputPathInputModel())
	input, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	inputModel, ok := input.(OutputPathInputText)
	if ok {
		utils.SetOutputDir(inputModel.value)
	}
}

func outputPathInputModel() OutputPathInputText {
	ti := textinput.New()
	ti.Placeholder = "Output Path Directory"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
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

func (m OutputPathInputText) View() string {
	return fmt.Sprintf(
		"\nEnter the relative output path(use `pwd` )?\n\n%s",
		m.textInput.View(),
	) + "\n"
}
