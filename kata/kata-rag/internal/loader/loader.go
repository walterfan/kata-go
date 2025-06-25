package loader

import (
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

type File struct {
    Path    string
    Content string
}

var CodeFileSuffixes []string = []string{".go", ".java", ".py", ".h", ".c", ".hpp", ".cpp"}

func LoadCodeFiles(dir string) ([]File, error) {
    var files []File
    err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
		for _, suffix := range CodeFileSuffixes {
			if strings.HasSuffix(path, suffix) {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				files = append(files, File{
					Path:    path,
					Content: string(content),
				})
			}
		}
        
        return nil
    })
    return files, err
}
