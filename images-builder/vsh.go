package main

import (
	"bytes"
	"fmt"
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
Screenshot {{ .Out }}.png
Sleep 500ms
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
	var a artefactReport
	a.Theme = theme
	a.Language = language
	a.SourceFile = sourceFilePath
	a.OutFile = outputFilePath + ".png"
	key := fmt.Sprintf("%s-%s-%s", theme, language, filepath.Base(sourceFilePath))
	f, err := os.CreateTemp("", key+".tape")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tpl.Execute(f, map[string]any{
		"File":  sourceFilePath,
		"Theme": theme,
		"Out":   outputFilePath,
	})
	err = f.Sync()
	if err != nil {
		a.Error = err
		return a
	}
	cmd := "vhs " + f.Name() + ""
	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput()
	_ = output // Uncomment to see the output

	fmt.Printf("%s", output)
	if err != nil {
		a.Error = fmt.Errorf("error running command %s: %w", cmd, err)
		return a
	}
	a.Success = true
	return a
}
