package config

import "path/filepath"

// Config はプロジェクト設定を表す
type Config struct {
	SourceRoot string   `yaml:"source_root"`
	PactRoot   string   `yaml:"pact_root"`
	OutputDir  string   `yaml:"output_dir"`
	Language   string   `yaml:"language"`
	Diagrams   []string `yaml:"diagrams"`
	Exclude    []string `yaml:"exclude"`
}

// Default はデフォルト設定を返す
func Default() *Config {
	return &Config{
		SourceRoot: "./src",
		PactRoot:   "./.pact",
		OutputDir:  "./diagrams",
		Language:   "go",
		Diagrams:   []string{"class", "sequence", "state", "flow"},
		Exclude:    []string{},
	}
}

// DiagramEnabled は指定した図が有効かどうかを返す
func (c *Config) DiagramEnabled(diagramType string) bool {
	for _, d := range c.Diagrams {
		if d == "all" || d == diagramType {
			return true
		}
	}
	return false
}

// IsExcluded は指定したパスが除外対象かどうかを返す
func (c *Config) IsExcluded(path string) bool {
	for _, pattern := range c.Exclude {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
		// ディレクトリ名でのマッチも試みる
		matched, err = filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}
	}
	return false
}
