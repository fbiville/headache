package main

import (
	"bufio"
	"fmt"
	tpl "html/template"
	"io"
	"os"
	"regexp"
	"strings"
)

type Configuration struct {
	HeaderFile   string            `json:"headerFile"`
	CommentStyle string            `json:"style"`
	Includes     []string          `json:"includes"`
	Excludes     []string          `json:"excludes"`
	TemplateData map[string]string `json:"data"`
}

type configuration struct {
	HeaderContents string
	HeaderRegex    *regexp.Regexp
	Includes       []string
	Excludes       []string
	writer         io.Writer
}

func ParseConfiguration(config Configuration) (*configuration, error) {
	contents, err := parseTemplate(config.HeaderFile, config.TemplateData, newCommentStyle(config.CommentStyle))
	if err != nil {
		return nil, err
	}
	return &configuration{
		HeaderContents: contents.actualContent,
		HeaderRegex:    contents.detectionRegex,
		Includes:       config.Includes,
		Excludes:       config.Excludes,
	}, nil
}

type templateResult struct {
	actualContent  string
	detectionRegex *regexp.Regexp
}

func parseTemplate(file string, data map[string]string, style CommentStyle) (*templateResult, error) {
	rawLines, err := readLines(file)
	if err != nil {
		return nil, err
	}
	commentedLines, err := applyComments(rawLines, style)
	if err != nil {
		return nil, err
	}
	template, err := tpl.New("header").Parse(strings.Join(commentedLines, "\n"))
	if err != nil {
		return nil, err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, data)
	if err != nil {
		return nil, err
	}
	regex, err := computeDetectionRegex(rawLines, data)
	if err != nil {
		return nil, err
	}
	return &templateResult{
		actualContent:  builder.String(),
		detectionRegex: regexp.MustCompile(regex),
	}, nil
}

func readLines(file string) ([]string, error) {
	openFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer openFile.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func applyComments(lines []string, style CommentStyle) ([]string, error) {
	result := make([]string, 0)
	if style.opening() {
		result = append(result, style.open())
	}
	for _, line := range lines {
		result = append(result, style.apply(line))
	}
	if style.closing() {
		result = append(result, style.close())
	}
	return result, nil
}

func computeDetectionRegex(lines []string, data map[string]string) (string, error) {
	regex := regexLines(lines)
	return injectDataRegex(strings.Join(regex, ""), data)
}

func injectDataRegex(result string, data map[string]string) (string, error) {
	template, err := tpl.New("header-regex").Parse(result)
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, regexValues(&data))
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func regexLines(lines []string) []string {
	result := make([]string, 0)
	result = append(result, "(?m)")
	result = append(result, "(?:\\/\\*\n)?")
	for _, line := range lines {
		result = append(result, fmt.Sprintf("%s\\Q%s\\E\n?", "(?:\\/{2}| \\*) ?", line))
	}
	result = append(result, "(?: \\*\\/)?")
	return result
}

func regexValues(data *map[string]string) *map[string]string {
	for k := range *data {
		(*data)[k] = "\\E.*\\Q"
	}
	return data
}
