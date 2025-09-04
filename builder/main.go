package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func main() {
	snippetFolderPathPtr := flag.String("snippet-folder-path", "../ressources/snippets-per-language", "the location of the snippet folder ")
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

	r, err := GetSnippetResource(fs, "../ressources/snippets-per-language")
	if err != nil {
		slog.Error("failed to read snippets ressources", slog.String("error", err.Error()))
	}
	slog.Info("successfully read snippets ressources", slog.Int("nb_snippets", r.snippetNumber()), slog.Int("nb_themes", len(THEMES)))
	for _, theme := range THEMES {
		for lang, snippets := range r.langages {
			for _, s := range snippets {

				snippetFilePath := filepath.Join(snippetFolderPath, lang, s)
				outputFilePath := filepath.Join(outputFolderPath, theme, lang, s) + ".png"
				err = fs.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm)
				if err != nil {
					slog.Error("failed to make out dirs", slog.String("error", err.Error()))
				}

				err = makeScreenshot(theme, lang, snippetFilePath, outputFilePath)
				if err != nil {
					slog.Error("make screenshot failed", slog.String("error", err.Error()), slog.String("snippetFilePath", snippetFilePath), slog.String("outputFilePath", outputFilePath))
				} else {
					slog.Info("make screenshot successful", slog.String("snippetFilePath", snippetFilePath), slog.String("outputFilePath", outputFilePath))
				}

			}
		}
	}
}
