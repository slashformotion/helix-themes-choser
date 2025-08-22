package main

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type SnippetResources struct {
	langages map[string][]string
}

func GetSnippetResource(fs afero.Fs, SnippetResourcesFolderPath string) (SnippetResources, error) {
	resp := SnippetResources{
		langages: make(map[string][]string, 0),
	}
	folderExist, err := afero.DirExists(fs, SnippetResourcesFolderPath)
	if err != nil {
		return SnippetResources{}, err
	}

	if !folderExist {
		panic("sdsd ")
	}

	languagesFolders, err := afero.ReadDir(fs, SnippetResourcesFolderPath)
	if err != nil {
		return SnippetResources{}, err
	}

	for _, ls := range languagesFolders {
		if !ls.IsDir() {
			continue
		}
		languageString := ls.Name()
		snippetFolderFiles, err := afero.ReadDir(fs, filepath.Join(SnippetResourcesFolderPath, languageString))
		if err != nil {
			return SnippetResources{}, err
		}
		if _, ok := resp.langages[languageString]; !ok {
			resp.langages[languageString] = make([]string, 0)
		}

		for _, s := range snippetFolderFiles {
			if s.IsDir() {
				continue
			}
			snippetName := s.Name()
			resp.langages[languageString] = append(
				resp.langages[languageString],
				snippetName,
			)

		}
	}

	return resp, nil
}
