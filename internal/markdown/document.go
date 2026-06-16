package markdown

type Document struct {
	Path        string
	Dir         string
	Source      string
	Body        string
	BodyLine    int
	FrontMatter map[string]string
	Headings    []Heading
	Links       []Link
	Images      []Image
}

type Heading struct {
	Level int
	Text  string
	Line  int
}

type Link struct {
	Text string
	URL  string
	Line int
}

type Image struct {
	Alt  string
	Path string
	Line int
}
