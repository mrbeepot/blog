package main

import (
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

func main() {
	filepath.WalkDir("./content", processWalk)
}

func processWalk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Fatal("Could not walk directory", err)
	}
	if !d.IsDir() {
		// A file
		info, infoErr := d.Info()
		if infoErr != nil {
			log.Fatal("Could not fetch info", infoErr)
		}
		buildHTML(path, info)
	}
	return nil
}

func buildHTML(path string, fileInfo fs.FileInfo) {
	err := writeHTML(parseMarkdown(path, fileInfo), path, fileInfo)
	println("path = ", path)
	if err != nil {
		log.Fatal("Could not write HTML file", err)
	}
}

func parseMarkdown(path string, fileInfo fs.FileInfo) []byte {
	// fullPath := filepath.Join(path, fileInfo.Name())
	parser := getMarkdownParser()
	renderer := getHTMLRenderer()
	mdFile, readFileError := ioutil.ReadFile(path)
	if readFileError != nil {
		log.Fatal("could not read", path)
	}
	return markdown.ToHTML(mdFile, parser, renderer)
}

func writeHTML(html []byte, path string, fileInfo fs.FileInfo) error {
	fullPath := filepath.Join(getOutDir(), getHTMLFileName(path))
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, os.ModePerm)
	return ioutil.WriteFile(fullPath, html, 0666)
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
