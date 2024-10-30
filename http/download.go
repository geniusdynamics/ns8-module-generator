package http

import (
	"fmt"
	"io"
	"net/http"
	"ns8-module-generator/processors"
	"ns8-module-generator/utils"
	"os"
)

type WriteCounter struct {
	Total uint64
}

func DownloadTemplate() {
	resp, err := http.Get(utils.TemplateZipURL)
	if err != nil {
		fmt.Printf("An error occurred while downloading template: %v", err)
	}
	defer resp.Body.Close()

	tempOut, err := os.Create("templatezip1.tmp")
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
	}

	_, err = io.Copy(tempOut, resp.Body)
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
	}
	err = os.Rename("templatezip1.tmp", "templatezip.zip")
	if err != nil {
		fmt.Printf("An err occurred: %v", err)
	}

	// Unzip Files After that
	processors.UnzipFiles("template", "template.zip")
}
