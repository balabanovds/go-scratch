package app

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"github.com/balabanovds/go-scratch/internal/core"
)

type App struct {
	outDir string
}

func New(outDir string) *App {
	return &App{
		outDir: outDir,
	}
}

func (a *App) Generate(data []core.Data) error {
	for _, d := range data {
		if err := a.processTemplate(d.TemplateFile, d.OutFile, d.Data); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) processTemplate(filename string, outfile string, data any) error {
	tmpl := template.Must(template.New("").ParseFiles(filename))
	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, filename, data)
	if err != nil {
		return fmt.Errorf("unable to parse data into template: %w", err)
	}
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		return fmt.Errorf("could not format processed template: %w", err)
	}
	outputPath := filepath.Join(a.outDir, outfile)
	fmt.Println("Writing file: ", outputPath)
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", outputPath, err)
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(formatted))
	if err != nil {
		return fmt.Errorf("failed to write string: %w", err)
	}

	return w.Flush()
}
