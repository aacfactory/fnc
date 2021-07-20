package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func LoadProject(path string) (p *Project, err error) {

	mod, modErr := loadModuleFile(path)
	if modErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", modErr)
		return
	}

	filenameMap := make(map[string]string)

	filenamesErr := loadFilenames(path, filenameMap)
	if filenamesErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", filenamesErr)
		return
	}

	if len(filenameMap) == 0 {
		err = fmt.Errorf("fnc load project failed, no go files loaded")
		return
	}

	filenames := make([]string, 0, 1)
	for k, _ := range filenameMap {
		filenames = append(filenames, k)
	}

	config := loader.Config{
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
		TypeChecker: types.Config{
			Error: func(err error) {

			},
		},
	}

	config.CreateFromFilenames(path, filenames...)

	program, loadErr := config.Load()
	if loadErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", loadErr)
		return
	}

	p = &Project{
		Path:      path,
		Module:    mod,
		Program:   program,
		Filenames: filenames,
	}

	return
}

type Module struct {
	Name     string
	Version  string
	Requires []Require
}

type Require struct {
	Name    string
	Version string
}

type Project struct {
	Path      string
	Module    Module
	Program   *loader.Program
	Filenames []string
}

func loadModuleFile(path string) (mod Module, err error) {
	modFilePath := filepath.Join(path, "go.mod")
	content, readErr := ioutil.ReadFile(modFilePath)
	if readErr != nil {
		err = fmt.Errorf("read go.mod failed, %v", readErr)
		return
	}
	if content == nil || len(content) == 0 {
		err = fmt.Errorf("read go.mod failed, no content")
		return
	}

	buf := bufio.NewReader(bytes.NewReader(content))

	reqFlag := false
	for {
		lineContent, _, readLineErr := buf.ReadLine()
		if readLineErr != nil {
			if readLineErr == io.EOF {
				break
			}
			err = fmt.Errorf("reade go.mod failed, %v", readLineErr)
		}
		line := strings.TrimSpace(string(lineContent))
		if len(line) == 0 {
			continue
		}
		// name
		if idx := strings.Index(line, "module "); idx == 0 {
			mod.Name = strings.TrimSpace(line[7+idx:])
			continue
		}
		// version
		if idx := strings.Index(line, "go "); idx == 0 {
			mod.Version = strings.TrimSpace(line[3+idx:])
			continue
		}
		// require
		if strings.Index(line, "require (") == 0 {
			reqFlag = true
			continue
		}
		if reqFlag {
			if strings.Index(line, ")") == 0 {
				reqFlag = false
				continue
			}
			if strings.Index(line, "\t") == 0 {
				line = line[2:]
			}
			items := strings.Split(line, " ")
			mod.Requires = append(mod.Requires, Require{
				Name:    items[0],
				Version: items[1],
			})
		}
	}

	return
}

func loadFilenames(path0 string, filenames map[string]string) (err error) {

	walkErr := filepath.Walk(path0, func(path string, info fs.FileInfo, cause error) (err error) {
		if strings.Contains(path, "/.git/") {
			return
		}
		if info.IsDir() {
			if path == path0 || strings.HasSuffix(path, ".git") {
				return
			}
			subErr := loadFilenames(path, filenames)
			if subErr != nil {
				err = subErr
				return
			}
			return
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return
		}
		if _, has := filenames[path]; !has {
			filenames[path] = info.Name()
		}

		return
	})

	if walkErr != nil {
		err = fmt.Errorf("load filename from %s failed, %v", path0, walkErr)
		return
	}

	return
}
