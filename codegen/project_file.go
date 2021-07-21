package codegen

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func loadPackageFiles(root string, modName string, pfs map[string][]string) (err error) {

	if strings.Contains(root, "\\") {
		root = strings.ReplaceAll(root, "\\", "/")
	}

	if strings.Contains(root, "/.") {
		return
	}

	modIdx := strings.Index(root, modName)
	if modIdx < 0 {
		err = fmt.Errorf("load project files failed, module name is not in path")
		return
	}
	name := root[modIdx:]

	walkErr := filepath.WalkDir(root, func(path string, dir fs.DirEntry, cause error) (err error) {
		if root == path {
			return
		}
		if strings.Contains(path, "\\") {
			path = strings.ReplaceAll(path, "\\", "/")
		}
		if strings.Contains(path, "/.git/") {
			return
		}
		if dir.IsDir() {
			subErr := loadPackageFiles(path, modName, pfs)
			if subErr != nil {
				err = subErr
				return
			}
			return
		}

		info, infoErr := dir.Info()
		if infoErr != nil {
			err = infoErr
			return
		}
		if info.IsDir() {
			return
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return
		}

		subPath := path[strings.Index(path, name) + len(name) +1:]
		if strings.Index(subPath, "/") > 0 {
			return
		}

		files, hasPkg := pfs[name]
		if !hasPkg {
			files = make([]string, 0, 1)
		}

		has := false
		for _, filePath := range files {
			if filePath == path {
				has = true
				break
			}
		}
		if has {
			return
		}
		files = append(files, path)
		pfs[name] = files
		return
	})

	if walkErr != nil {
		err = fmt.Errorf("load filename from %s failed, %v", root, walkErr)
		return
	}

	return
}
