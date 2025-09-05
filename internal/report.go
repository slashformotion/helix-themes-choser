package internal

type Snippet struct {
	URL  string
	Name string
}

type Languages struct {
	Lang     string
	Snippets []Snippet
}

type Theme struct {
	Name      string
	Languages []Languages
}

type Report struct {
	Themes []Theme
}
