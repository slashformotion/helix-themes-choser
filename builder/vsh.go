package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const tplString = `
Set Padding 0
Set Width 1920
Set Height 1080
Set FontSize 30

Require hx
Type "hx {{ .File }} -c config.toml"
Enter
Wait+Screen@2s     /NOR/
Type ":theme {{ .Theme }}"
Sleep 10ms
Enter
Sleep 100ms
Screenshot {{ .Out }}
Sleep 50ms
Type ":q"
Sleep 100ms
`

var tpl = template.Must(template.New("main").Parse(tplString))

func renderVhsTemplateString(theme, snippetFilePath, screenshotFilePath string) ([]byte, error) {
	var b bytes.Buffer
	err := tpl.Execute(&b, map[string]any{
		"File":  snippetFilePath,
		"Theme": theme,
		"Out":   screenshotFilePath,
	})
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func makeScreenshot(theme, languageKey, snippetFilePath, outputFilePath string) error {
	f, err := os.CreateTemp("", filepath.Base(outputFilePath)+".tape")
	if err != nil {
		return fmt.Errorf("failed to create .tape file for %s %s %s: %w", theme, languageKey, snippetFilePath, err)
	}
	defer f.Close()
	tpl.Execute(f, map[string]any{
		"File":  snippetFilePath,
		"Theme": theme,
		"Out":   outputFilePath,
	})
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("failed to flush file to disk: %w", err)
	}
	cmd := "vhs " + f.Name() + ""
	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput()
	if err != nil {
		slog.Debug("print combined output", slog.String("output", string(output)))
		return fmt.Errorf("failed to run vhs: %w", err)
	}
	return nil
}
