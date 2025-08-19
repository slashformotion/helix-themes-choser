package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
)

const LANGUAGE_FOLDER = "./languages/"
const OUT_FOLDER = "./out/"
const defaultArtefactFile = "artefacts.json"

var languagesMap = map[string][]string
	"go":     {"main.go"},
	"python": {"main.py"},
}

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

type job struct {
	theme    string
	language string
	langFile string
}

type artefactReport struct {
	Success    bool   `json:"-"`
	Error      error  `json:"-"`
	Theme      string `json:"theme"`
	SourceFile string `json:"source_file"`
	Language   string `json:"language"`
	OutFile    string `json:"out_file"`
}

type report map[string]map[string][]string

func (a artefactReport) Id() string {
	return fmt.Sprintf("%s-%s-%s", a.Theme, a.Language, filepath.Base(a.SourceFile))
}

type artefactReports []artefactReport

var tpl = template.Must(template.New("main").Parse(tplString))

func main() {
	r := make(report, 0)
	topic := make(chan job, len(themes)*len(languagesMap)*10)
	artefactReportChan := make(chan artefactReport, len(themes)*len(languagesMap)*10)
	for _, theme := range themes {
		for language, langFiles := range languagesMap {
			for _, langFile := range langFiles {
				r.
				topic <- job{
					theme:    theme,
					language: language,
					langFile: langFile,
				}
			}
		}
	}
	close(topic)
	wg := sync.WaitGroup{}
	for range 1 {
		wg.Add(1)
		go func() {
			for j := range topic {
				artefactReportChan <- makeScreenshot(
					j.theme,
					j.language,
					filepath.Join(LANGUAGE_FOLDER, j.langFile),
					filepath.Join(OUT_FOLDER, fmt.Sprintf("%s-%s-%s", j.language, j.theme, j.langFile)),
				)
				fmt.Printf("%v", j)
			}

			wg.Done()
		}()
	}
	wg.Wait()
	close(artefactReportChan)
	artefacts := artefactReports{}
	for a := range artefactReportChan {
		if a.Error != nil {
			fmt.Printf("Error running %s: %v\n", a.Id(), a.Error)
		} else {
			artefacts = append(artefacts, a)
			fmt.Printf("Success: %s %s %s\n", a.Theme, a.Language, a.SourceFile)
		}
	}
	artfactFileContent, err := json.MarshalIndent(&artefacts, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling artefacts: %v\n", err)
		return
	}

	os.WriteFile(OUT_FOLDER+defaultArtefactFile, artfactFileContent, 0644)

}

func makeScreenshot(theme, language, sourceFilePath, outputFilePath string) artefactReport {
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

var themes = []string{
	"rasmus",
	// "vintage",
	// "starlight",
	// "gruber-darker",
	// "solarized_dark",
	// "github_light",
	// "ayu_evolve",
	// "material_deep_ocean",
	// "gruvbox_dark_soft",
	// "catppuccin_frappe",
	// "modus_vivendi",
	// "rose_pine_moon",
	// "github_light_colorblind",
	// "everblush",
	// "naysayer",
	// "hex_toxic",
	// "zenburn",
	// "eiffel",
	// "tokyonight_moon",
	// "catppuccin_macchiato",
	// "hex_steel",
	// "nord-night",
	// "github_dark_colorblind",
	// "base16_terminal",
	// "monokai_soda",
	// "seoul256-dark",
	// "seoul256-light-hard",
	// "term16_light",
	// "tokyonight",
	// "seoul256-dark-hard",
	// "iceberg-light",
	// "ferra",
	// "cyan_light",
	// "bogster",
	// "yo_light",
	// "base16_default_light",
	// "bogster_light",
	// "beans",
	// "iroaseta",
	// "gruvbox_light_hard",
	// "tokyonight_day",
	// "monokai_aqua",
	// "monokai_pro_machine",
	// "horizon-dark",
	// "molokai",
	// "monokai_pro",
	// "snazzy",
	// "kaolin-valley-dark",
	// "acme",
	// "flexoki_light",
	// "material_darker",
	// "hex_poison",
	// "darcula-solid",
	// "autumn",
	// "github_dark_tritanopia",
	// "ingrid",
	// "ayu_light",
	// "github_light_tritanopia",
	// "zed_onelight",
	// "catppuccin_latte",
	// "zed_onedark",
	// "tokyonight_storm",
	// "noctis",
	// "yellowed",
	// "base16_default_dark",
	// "monokai_pro_octagon",
	// "hex_lavender",
	// "jetbrains_dark",
	// "ttox",
	// "github_dark_dimmed",
	// "varua",
	// "rose_pine_dawn",
	// "modus_vivendi_tritanopia",
	// "papercolor-light",
	// "material_palenight",
	// "kaolin-dark",
	// "vim_dark_high_contrast",
	// "seoul256-light",
	// "heisenberg",
	// "fleet_dark",
	// "voxed",
	// "seoul256-light-soft",
	// "nightfox",
	// "noctis_bordo",
	// "merionette",
	// "nord_light",
	// "seoul256-dark-soft",
	// "solarized_light",
	// "serika-light",
	// "gruvbox_light_soft",
	// "modus_operandi_deuteranopia",
	// "github_dark_high_contrast",
	// "modus_vivendi_tinted",
	// "catppuccin_mocha",
	// "yo",
	// "ayu_mirage",
	// "monokai_pro_ristretto",
	// "autumn_night",
	// "everforest_light",
	// "material_oceanic",
	// "flexoki_dark",
	// "mellow",
	// "meliora",
	// "onedarker",
	// "flatwhite",
	// "kanagawa",
	// "gruvbox_dark_hard",
	// "github_light_high_contrast",
	// "modus_operandi_tritanopia",
	// "adwaita-dark",
	// "poimandres_storm",
	// "darcula",
	// "iceberg-dark",
	// "jellybeans",
	// "amberwood",
	// "everforest_dark",
	// "modus_vivendi_deuteranopia",
	// "serika-dark",
	// "curzon",
	// "onedark",
	// "carbonfox",
	// "poimandres",
	// "dark_plus",
	// "onelight",
	// "boo_berry",
	// "ayu_dark",
	// "dark_high_contrast",
	// "pop-dark",
	// "monokai_pro_spectrum",
	// "emacs",
	// "dracula",
	// "kaolin-light",
	// "modus_operandi",
	// "term16_dark",
	// "gruvbox",
	// "new_moon",
	// "adwaita-light",
	// "dracula_at_night",
	// "github_dark",
	// "ao",
	// "kanagawa-dragon",
	// "spacebones_light",
	// "papercolor-dark",
	// "doom_acario_dark",
	// "night_owl",
	// "modus_operandi_tinted",
	// "base16_transparent",
	// "monokai",
	// "nord",
	// "sonokai",
	// "yo_berry",
	// "gruvbox_light",
	// "rose_pine",
	// "penumbra+",
	// "sunset",
}
