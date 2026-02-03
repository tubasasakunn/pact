package main

import (
	"os"
	"path/filepath"
)

// expandFiles はファイルパターンのリストを展開し、実際のファイルパスのリストを返す
// ディレクトリが指定された場合は *.pact ファイルを自動展開する
func expandFiles(patterns []string) []string {
	var files []string
	for _, pattern := range patterns {
		info, err := os.Stat(pattern)
		if err == nil && info.IsDir() {
			matches, _ := filepath.Glob(filepath.Join(pattern, "*.pact"))
			files = append(files, matches...)
		} else {
			matches, _ := filepath.Glob(pattern)
			if len(matches) == 0 {
				files = append(files, pattern)
			} else {
				files = append(files, matches...)
			}
		}
	}
	return files
}
