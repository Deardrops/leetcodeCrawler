package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const problemBaseCN = "https://leetcode-cn.com/problems/"
const problemBase = "https://leetcode.com/problems/"

func createQuestion(basePath string, item *Item, question *Question) error {
	titleSlug := item.Question.TitleSlug
	folder := path.Join(basePath, titleSlug)
	err := createFolder(folder)
	if err != nil {
		return err
	}
	defaultCode := parseDefaultCode(question.CodeDefinition)
	err = createSolutionFile(folder, defaultCode)
	if err != nil {
		return err
	}
	err = createReadmeFile(folder, item.Title, titleSlug, question.TranslatedContent)
	if err != nil {
		return err
	}
	return nil
}

func createFolder(folder string) error {
	err := os.Mkdir(folder, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("folder existed: ", folder)
			return err
		}
	}
	return nil
}

func createSolutionFile(folder, defaultCode string) error {
	solutionPath := path.Join(folder, "solution.go")
	f, err := os.Create(solutionPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s\n", solutionPath)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "package main")
	fmt.Fprintln(w)
	fmt.Fprint(w, defaultCode)
	return w.Flush()
}

func createReadmeFile(folder, title, titleSlug, content string) error {
	readmePath := path.Join(folder, "readme.md")
	f, err := os.Create(readmePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s\n", readmePath)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "## "+title)
	fmt.Fprintln(w)
	fmt.Fprint(w, content)
	fmt.Fprintln(w)
	fmt.Fprintln(w,"-----")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "### 链接：")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "中文："+problemBaseCN+titleSlug)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "英文："+problemBase+titleSlug)
	return w.Flush()
}

func parseDefaultCode(codeDeinition string) string {
	var CodeDefinitions []CodeDefinition
	bytes := []byte(codeDeinition)
	err := json.Unmarshal(bytes, &CodeDefinitions)
	if err != nil {
		fmt.Println("Failed to parse code definition json: ", err)
		return ""
	}
	for _, item := range CodeDefinitions {
		if item.Value == "golang" {
			return item.DefaultCode
		}
	}
	return ""
}

type CodeDefinition struct {
	Value       string
	Text        string
	DefaultCode string
}
