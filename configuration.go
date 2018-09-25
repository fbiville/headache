package header

import (
	"bufio"
	tpl "html/template"
	"io"
	"os"
	"strings"
)

type configuration struct {
	HeaderContents string
	Includes       []string
	Excludes       []string
	writer         io.Writer
}

func NewConfiguration(headerFile string,
	style CommentStyle,
	includes []string,
	excludes []string,
	data map[string]string) (*configuration, error) {
	contents, err := parseTemplate(headerFile, data, style)
	if err != nil {
		return nil, err
	}
	return &configuration{
		HeaderContents: contents,
		Includes:       includes,
		Excludes:       excludes,
	}, nil
}

func parseTemplate(file string, data map[string]string, style CommentStyle) (string, error) {
	lines, err := readStyledLines(file, style)
	if err != nil {
		return "", err
	}
	template, err := tpl.New("header").Parse(strings.Join(lines, "\n"))
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func readStyledLines(file string, style CommentStyle) ([]string, error) {
	openFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer openFile.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(openFile)
	if style.opening() {
		lines = append(lines, style.open())
	}
	for scanner.Scan() {
		lines = append(lines, style.apply(scanner.Text()))
	}
	if style.closing() {
		lines = append(lines, style.close())
	}
	return lines, nil
}
