package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/afero"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	logger.InfoContext(ctx, "starting")
	fs := afero.NewOsFs()

	r, err := GetSnippetResource(fs, "../ressources/snippets-per-language")
	if err != nil {
		panic(err)
	}
	logger.Info("sdsd", slog.Any("sds", r))

	// r := make(report, 0)
	// topic := make(chan job, len(themes)*len(languagesMap)*10)
	// artefactReportChan := make(chan artefactReport, len(themes)*len(languagesMap)*10)
	// for _, theme := range themes {
	// 	for language, langFiles := range languagesMap {
	// 		for _, langFile := range langFiles {
	// 			r.
	// 				topic <- job{
	// 				theme:    theme,
	// 				language: language,
	// 				langFile: langFile,
	// 			}
	// 		}
	// 	}
	// }
	// close(topic)
	// wg := sync.WaitGroup{}
	// for range 1 {
	// 	wg.Add(1)
	// 	go func() {
	// 		for j := range topic {
	// 			artefactReportChan <- makeScreenshot(
	// 				j.theme,
	// 				j.language,
	// 				filepath.Join(LANGUAGE_FOLDER, j.langFile),
	// 				filepath.Join(OUT_FOLDER, fmt.Sprintf("%s-%s-%s", j.language, j.theme, j.langFile)),
	// 			)
	// 			fmt.Printf("%v", j)
	// 		}

	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()
	// close(artefactReportChan)
	// artefacts := artefactReports{}
	// for a := range artefactReportChan {
	// 	if a.Error != nil {
	// 		fmt.Printf("Error running %s: %v\n", a.Id(), a.Error)
	// 	} else {
	// 		artefacts = append(artefacts, a)
	// 		fmt.Printf("Success: %s %s %s\n", a.Theme, a.Language, a.SourceFile)
	// 	}
	// }
	// artfactFileContent, err := json.MarshalIndent(&artefacts, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Error marshalling artefacts: %v\n", err)
	// 	return
	// }

	// os.WriteFile(OUT_FOLDER+defaultArtefactFile, artfactFileContent, 0644)

}
