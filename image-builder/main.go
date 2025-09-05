package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/slashformotion/helix-theme-choser/internal"
)

func main() {
	snippetFolderPathPtr := flag.String("snippet-folder-path", "ressources/snippets-per-language", "the location of the snippet folder ")
	outputFolderPathPtr := flag.String("output-folder-path", "out", "the location of the output folder ")
	flag.Parse()
	snippetFolderPath := *snippetFolderPathPtr
	outputFolderPath := *outputFolderPathPtr
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)
	slog.InfoContext(ctx, "starting")
	fs := afero.NewOsFs()

	r, err := internal.GetSnippetResource(fs, snippetFolderPath)
	if err != nil {
		slog.Error("failed to read snippets ressources", slog.String("error", err.Error()))
	}
	slog.Info("successfully read snippets ressources", slog.Int("nb_snippets", r.SnippetNumber()), slog.Int("nb_themes", len(internal.THEMES)))
	report := internal.Report{}
	for _, theme := range internal.THEMES {
		themeReport := internal.Theme{Name: theme}
		for lang, snippets := range r.Langages {
			langReport := internal.Languages{
				Lang: lang,
			}
			for _, snippet := range snippets {
				snippetFilePath := filepath.Join(snippetFolderPath, lang, snippet)
				outputFilePath := filepath.Join(outputFolderPath, theme, lang, snippet) + ".png"
				err = fs.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm)
				if err != nil {
					slog.Error("failed to make out dirs", slog.String("error", err.Error()))
				}

				err = internal.MakeScreenshot(theme, lang, snippetFilePath, outputFilePath)
				if err != nil {
					slog.Error("make screenshot failed", slog.String("error", err.Error()), slog.String("snippetFilePath", snippetFilePath), slog.String("outputFilePath", outputFilePath))
				} else {
					snippetReport := internal.Snippet{Name: snippet, URL: fmt.Sprintf("./%s/%s/%s.png", theme, lang, snippet)}
					langReport.Snippets = append(langReport.Snippets, snippetReport)
					slog.Info("make screenshot successful", slog.String("snippetFilePath", snippetFilePath), slog.String("outputFilePath", outputFilePath))
				}
			}
			themeReport.Languages = append(themeReport.Languages, langReport)
		}
		report.Themes = append(report.Themes, themeReport)
	}

	reportBytes, err := json.MarshalIndent(&report, "", "  ")
	if err != nil {
		panic(err)
	}

	err = afero.WriteFile(fs, filepath.Join(outputFolderPath, "report.json"), reportBytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

}
