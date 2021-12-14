package main

import (
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type BlogPage struct {
	Title   string
	Css     string
	Js      string
	Content template.HTML
}

func main() {
	filepath.WalkDir("./content", processWalk)
}

func processWalk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Fatal("Could not walk directory", err)
	}
	if !d.IsDir() {
		info, infoErr := d.Info()
		if infoErr != nil {
			log.Fatal("Could not fetch info", infoErr)
		}
		buildHTML(path, info)
	}
	return nil
}

func buildHTML(path string, fileInfo fs.FileInfo) {
	err := writeHTMLTemplate(parseMarkdown(path, fileInfo), path, fileInfo)
	println("path = ", path)
	if err != nil {
		log.Fatal("Could not write HTML file", err)
	}
}

func parseMarkdown(path string, fileInfo fs.FileInfo) []byte {
	parser := getMarkdownParser()
	renderer := getHTMLRenderer()
	mdFile, readFileError := ioutil.ReadFile(path)
	if readFileError != nil {
		log.Fatal("could not read", path)
	}
	return markdown.ToHTML(mdFile, parser, renderer)
}

// func writeHTML(html []byte, path string, fileInfo fs.FileInfo) error {
// 	fullPath := filepath.Join(getOutDir(), getHTMLFileName(path))
// 	dir := filepath.Dir(fullPath)
// 	os.MkdirAll(dir, os.ModePerm)
// 	return ioutil.WriteFile(fullPath, html, 0666)
// }

func writeHTMLTemplate(htmlContent []byte, path string, fileInfo fs.FileInfo) error {
	blogPage := BlogPage{
		Title:   strings.ReplaceAll(filepath.Base(path), ".md", ""),
		Css:     getCss(),
		Js:      getJs(),
		Content: template.HTML(htmlContent),
	}
	index, _ := ioutil.ReadFile("./templates/index.html")
	temp, tempErr := template.New("blogPage").Parse(string(index))
	if tempErr != nil {
		log.Fatal("Could not parse ", tempErr)
	}
	fullPath := filepath.Join(getOutDir(), getHTMLFileName(path))
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, os.ModePerm)
	fileWriter, _ := os.Create(fullPath)
	return temp.Execute(fileWriter, blogPage)

}

func getCss() string {
	return ""
}

func getJs() string {
	return ""
}

func getHTMLFileName(fileInfo string) string {
	return modifyExtensionToHTML(removeSpaces(fileInfo))
}

func modifyExtensionToHTML(name string) string {
	return strings.ReplaceAll(name, ".md", ".html")
}

func removeSpaces(name string) string {
	return strings.ReplaceAll(name, " ", "")
}

func getMarkdownExtensions() parser.Extensions {
	return parser.CommonExtensions | parser.AutoHeadingIDs
}

func getMarkdownParser() *parser.Parser {
	return parser.NewWithExtensions(getMarkdownExtensions())
}

func getHTMLRenderer() *html.Renderer {
	return html.NewRenderer(getHTMLRendererOptions())
}

func getHTMLFlags() html.Flags {
	return html.CommonFlags | html.HrefTargetBlank
}

func getHTMLRendererOptions() html.RendererOptions {
	return html.RendererOptions{Flags: getHTMLFlags()}
}

func getOutDir() string {
	return "./parsed"
}
