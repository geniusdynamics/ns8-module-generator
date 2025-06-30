package http

import (
	"fmt"
	"io"
	"net/http"
	"ns8-module-generator/config"
	"ns8-module-generator/processors"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type WriteCounter struct {
	Total uint64
}

type ProgressModel struct {
	progress   progress.Model
	totalBytes int64
	readBytes  int64
	err        error
}
type (
	progressMsg float64
	errMsg      error
)

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case progressMsg:
		m.readBytes = int64(float64(m.totalBytes) * float64(msg))
		m.progress.SetPercent(float64(m.readBytes) / float64(m.totalBytes))
		if m.readBytes >= m.totalBytes {
			return m, tea.Quit
		}
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	return m, nil
}

// View renders the progress bar
func (m ProgressModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("An error occurred: %v\n", m.err)
	}
	return fmt.Sprintf("Downloading template...\n%s", m.progress.View())
}

func DownloadTemplate() error {
	// Fetch the file size
	resp, err := http.Head(config.Cfg.TemplateZipURL)
	if err != nil {
		return fmt.Errorf("Failed to get file size: %v", err)
	}
	defer resp.Body.Close()

	// Create the Bubble Tea program for progress
	fileSize := resp.ContentLength
	m := ProgressModel{
		progress:   progress.New(progress.WithDefaultGradient()),
		totalBytes: fileSize,
	}
	p := tea.NewProgram(m)

	errChan := make(chan error, 1)

	// Start downloading and show progress
	go func() {
		defer close(errChan)
		resp, err := http.Get(config.Cfg.TemplateZipURL)
		if err != nil {
			p.Send(errMsg(err))
			errChan <- fmt.Errorf("Failed to download template: %v", err)
			return
		}
		defer resp.Body.Close()

		outFile, err := os.Create("templatezip1.tmp")
		if err != nil {
			p.Send(errMsg(err))
			errChan <- fmt.Errorf("Failed to create temporary file: %v", err)
			return
		}
		defer outFile.Close()

		// Track download progress
		progressReader := &ProgressReader{Reader: resp.Body, TotalSize: fileSize, Program: p}
		_, err = io.Copy(outFile, progressReader)
		if err != nil {
			p.Send(errMsg(err))
			errChan <- fmt.Errorf("Failed to copy file content: %v", err)
			return
		}

		// Rename to final zip file
		if err := os.Rename("templatezip1.tmp", "templatezip.zip"); err != nil {
			p.Send(errMsg(err))
			errChan <- fmt.Errorf("Failed to rename temporary file: %v", err)
			return
		}

		// Unzip the files and update directory
		processors.UnzipFiles("template", "templatezip.zip")
		p.Send(progressMsg(1.0)) // Mark as complete
		errChan <- nil
	}()

	// Start the Bubble Tea program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error running progress UI: %v", err)
	}

	return <-errChan
}

// ProgressReader implements io.Reader with Bubble Tea progress updates
type ProgressReader struct {
	io.Reader
	TotalSize int64
	ReadSize  int64
	Program   *tea.Program
}

// Read implements the io.Reader interface with progress updates
func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.ReadSize += int64(n)

	progress := float64(r.ReadSize) / float64(r.TotalSize)
	r.Program.Send(progressMsg(progress))
	return n, err
}
