package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/slashformotion/helix-theme-choser/internal"
)

type Job struct {
	Theme       string
	Lang        string
	Snippet     string
	SnippetPath string
	OutputPath  string
}

type Result struct {
	Theme   string
	Lang    string
	Snippet internal.Snippet
	Err     error
}

func worker(id int, fs afero.Fs, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		err := fs.MkdirAll(filepath.Dir(job.OutputPath), os.ModePerm)
		if err != nil {
			slog.Error("failed to create output directories", slog.String("error", err.Error()))
			results <- Result{Theme: job.Theme, Lang: job.Lang, Err: err}
			continue
		}

		err = internal.MakeScreenshot(job.Theme, job.Lang, job.SnippetPath, job.OutputPath)
		if err != nil {
			slog.Error("make screenshot failed",
				slog.String("error", err.Error()),
				slog.String("snippetFilePath", job.SnippetPath),
				slog.String("outputFilePath", job.OutputPath),
			)
		} else {
			slog.Info("make screenshot successful",
				slog.String("snippetFilePath", job.SnippetPath),
				slog.String("outputFilePath", job.OutputPath),
			)
		}

		results <- Result{
			Theme:   job.Theme,
			Lang:    job.Lang,
			Snippet: internal.Snippet{Name: job.Snippet, URL: fmt.Sprintf("./%s/%s/%s.png", job.Theme, job.Lang, job.Snippet)},
			Err:     err,
		}
	}
}

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
		return
	}

	slog.Info("successfully read snippets ressources",
		slog.Int("nb_snippets", r.SnippetNumber()),
		slog.Int("nb_themes", len(internal.THEMES)),
	)

	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var wg sync.WaitGroup

	// Start 4 worker goroutines
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, fs, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for _, theme := range internal.THEMES {
			for lang, snippets := range r.Langages {
				for _, snippet := range snippets {
					snippetPath := filepath.Join(snippetFolderPath, lang, snippet)
					outputPath := filepath.Join(outputFolderPath, theme, lang, snippet) + ".png"
					jobs <- Job{
						Theme:       theme,
						Lang:        lang,
						Snippet:     snippet,
						SnippetPath: snippetPath,
						OutputPath:  outputPath,
					}
				}
			}
		}
		close(jobs)
	}()

	// Collect results
	report := internal.Report{}
	reportMap := make(map[string]map[string][]internal.Snippet)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			continue
		}
		if _, ok := reportMap[result.Theme]; !ok {
			reportMap[result.Theme] = make(map[string][]internal.Snippet)
		}
		reportMap[result.Theme][result.Lang] = append(reportMap[result.Theme][result.Lang], result.Snippet)
	}

	// Build report from map
	for theme, langs := range reportMap {
		themeReport := internal.Theme{Name: theme}
		for lang, snippets := range langs {
			langReport := internal.Languages{
				Lang:     lang,
				Snippets: snippets,
			}
			themeReport.Languages = append(themeReport.Languages, langReport)
		}
		report.Themes = append(report.Themes, themeReport)
	}

	// Write report.json
	reportBytes, err := json.MarshalIndent(&report, "", "  ")
	if err != nil {
		panic(err)
	}

	err = afero.WriteFile(fs, filepath.Join(outputFolderPath, "report.json"), reportBytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
