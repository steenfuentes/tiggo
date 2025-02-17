package analyze

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/steenfuentes/gogit/internal/git"
	"github.com/steenfuentes/gogit/internal/llm"
)

var SUM_PROMPT = `Generate a summary of the following code containing file. The summary
 should be an overview of the functionality that the code within the file offers. Do not
 include any introduction or sign off of your own in the response. Only include the summary in your response.`

var DIFF_SUM_PROMPT = `Use the supplied output of a 'git diff' command for a single file
  to summarize the changes made. Be very detailed and objective with the analyses, but
  make no attempts to infer the impact of the changes. Generate a change summary for each
  change in the diff and be sure to reference the line numbers within each summary. Do not include any introduction or
  sign off of your own in the response. Only include the summary in your response.`

func GetFileContent(path string) (*string, error) {

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error getting file: %w", err)
	}

	out_str := strings.TrimSpace(string(b))

	return &out_str, nil
}

type FileAnalysis struct {
	Path        string
	CodeSummary string
	DiffSummary string
}
type Analyzer struct {
	LLMClient    *llm.LLMClient
	DiffRange    *git.DiffRange
	FileAnalyses []FileAnalysis
}

func NewFileAnalysis(filePath string) *FileAnalysis {
	return &FileAnalysis{Path: filePath}
}

func (fa *FileAnalysis) AddFileSummary(llmClient *llm.LLMClient) error {

	fileContent, err := GetFileContent(fa.Path)
	if err != nil {
		return err
	}

	content := SUM_PROMPT + *fileContent

	resp, err := llmClient.SendMessage(content)
	if err != nil {
		return err
	}

	fa.CodeSummary = resp.Content[0].Text

	return nil

}

func (fa *FileAnalysis) AddDiffSummary(
	diffRange *git.DiffRange, llmClient *llm.LLMClient) error {

	fileDiff, err := diffRange.GetFileDiff(fa.Path)
	if err != nil {
		return err
	}

	content := DIFF_SUM_PROMPT + *fileDiff

	resp, err := llmClient.SendMessage(content)
	if err != nil {
		return err
	}

	fa.DiffSummary = resp.Content[0].Text

	return nil
}

func (fa *FileAnalysis) String() string {
	return fmt.Sprintf("File:\n%s\n\nCode Summary:\n%s\n\nDiff Summary:\n%s\n\n", fa.Path, fa.CodeSummary, fa.DiffSummary)
}

func (a *Analyzer) DoAnalysis() error {

	files, err := a.DiffRange.GetChangedFiles()
	if err != nil {
		return err
	}

	errChan := make(chan error, len(files))

	var wg sync.WaitGroup

	a.FileAnalyses = make([]FileAnalysis, len(files))

	sem := make(chan struct{}, 5)

	for i, file := range files {
		wg.Add(1)
		go func(i int, file string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			analysis := FileAnalysis{Path: file}

			sumChan := make(chan error, 1)
			diffChan := make(chan error, 1)

			go func() {
				sumChan <- analysis.AddFileSummary(a.LLMClient)
			}()

			go func() {
				diffChan <- analysis.AddDiffSummary(a.DiffRange, a.LLMClient)
			}()

			if err := <-sumChan; err != nil {
				errChan <- fmt.Errorf("error adding file summary for %s: %w", file, err)
				return
			}

			if err := <-diffChan; err != nil {
				errChan <- fmt.Errorf("error adding diff summary for %s: %w", file, err)
				return
			}

			a.FileAnalyses[i] = analysis
		}(i, file)

	}

	wg.Wait()

	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Analyzer) AnalysesAsString() *string {

	var output string

	for _, analysis := range a.FileAnalyses {
		output = output + analysis.String()
	}

	return &output
}

type PromptData struct {
	FileAnalyses string
}

func injectFileAnalyses(templatePath string, analyses string) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	data := PromptData{
		FileAnalyses: analyses,
	}

	tmpl, err := template.New("prompt").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (a *Analyzer) Summarize(outPath string) error {

	analyses := a.AnalysesAsString()

	prompt, err := injectFileAnalyses("system_prompt.md", *analyses)
	if err != nil {
		return err
	}

	resp, err := a.LLMClient.SendMessage(prompt)
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, []byte(resp.Content[0].Text), 0644)
	if err != nil {
		return err
	}

	return nil

}
