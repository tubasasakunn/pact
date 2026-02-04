package resolver

import (
	"os"
	"path/filepath"

	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
	"pact/internal/infrastructure/parser"
)

// Parser はファイルをパースするインターフェース
type Parser interface {
	ParseFile(path string) (*ast.SpecFile, error)
}

// Resolver はインポートを解決する
type Resolver struct {
	parser Parser
}

// NewResolver は新しいResolverを作成する
func NewResolver(parser Parser) *Resolver {
	return &Resolver{parser: parser}
}

// ImportResolver はパスベースのインポート解決器
type ImportResolver struct{}

// NewImportResolver は新しいImportResolverを作成する
func NewImportResolver() *ImportResolver {
	return &ImportResolver{}
}

// Resolve はファイルパスからインポートを解決し、依存順に返す
func (r *ImportResolver) Resolve(path string) ([]string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, &errors.ImportError{Path: path, Message: err.Error()}
	}

	visited := make(map[string]bool)
	inProgress := make(map[string]bool)
	var result []string

	if err := r.resolveFile(absPath, visited, inProgress, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// resolveFile は単一ファイルを解決する
func (r *ImportResolver) resolveFile(
	absPath string,
	visited map[string]bool,
	inProgress map[string]bool,
	result *[]string,
) error {
	// サイクル検出
	if inProgress[absPath] {
		return &errors.CycleError{Cycle: []string{absPath}}
	}

	// 既に処理済み
	if visited[absPath] {
		return nil
	}

	inProgress[absPath] = true
	defer func() { delete(inProgress, absPath) }()

	// ファイルを読み込んでパース
	content, err := os.ReadFile(absPath)
	if err != nil {
		return &errors.ImportError{Path: absPath, Message: err.Error()}
	}

	spec, err := parser.ParseString(string(content))
	if err != nil {
		return &errors.ImportError{Path: absPath, Message: "parse error", Cause: err}
	}

	// インポートを処理
	fileDir := filepath.Dir(absPath)
	for _, imp := range spec.Imports {
		impAbsPath, err := filepath.Abs(filepath.Join(fileDir, imp.Path))
		if err != nil {
			return &errors.ImportError{Pos: imp.Pos, Path: imp.Path, Message: err.Error()}
		}

		if err := r.resolveFile(impAbsPath, visited, inProgress, result); err != nil {
			return err
		}
	}

	visited[absPath] = true
	*result = append(*result, absPath)
	return nil
}

// IsCycleError はエラーがCycleErrorかどうかを判定する
func IsCycleError(err error) bool {
	_, ok := err.(*errors.CycleError)
	return ok
}

// Resolve はファイルリストのインポートを解決し、依存順にソートして返す
func (r *Resolver) Resolve(files []*ast.SpecFile) ([]*ast.SpecFile, error) {
	if len(files) == 0 {
		return []*ast.SpecFile{}, nil
	}

	// 全ファイルを一つのグラフとして処理
	visited := make(map[string]bool)
	inProgress := make(map[string]bool)
	fileMap := make(map[string]*ast.SpecFile)
	var result []*ast.SpecFile

	// ファイルマップを作成
	for _, f := range files {
		fileMap[f.Path] = f
	}

	// 各ファイルを処理
	for _, f := range files {
		if err := r.resolveFromMap(f, fileMap, visited, inProgress, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// resolveFromMap はファイルマップからインポートを解決する
func (r *Resolver) resolveFromMap(
	file *ast.SpecFile,
	fileMap map[string]*ast.SpecFile,
	visited map[string]bool,
	inProgress map[string]bool,
	result *[]*ast.SpecFile,
) error {
	path := file.Path

	// サイクル検出
	if inProgress[path] {
		return &errors.CycleError{Cycle: []string{path}}
	}

	// 既に処理済み
	if visited[path] {
		return nil
	}

	inProgress[path] = true
	defer func() { delete(inProgress, path) }()

	// インポートを処理
	for _, imp := range file.Imports {
		impPath := normalizePath(filepath.Dir(file.Path), imp.Path)
		impFile, ok := fileMap[impPath]
		if !ok {
			return &errors.ImportError{
				Pos:     imp.Pos,
				Path:    imp.Path,
				Message: "file not found",
			}
		}

		if err := r.resolveFromMap(impFile, fileMap, visited, inProgress, result); err != nil {
			return err
		}
	}

	visited[path] = true
	*result = append(*result, file)
	return nil
}

// ResolveFile は単一ファイルのインポートを解決する
func (r *Resolver) ResolveFile(file *ast.SpecFile, basePath string) ([]*ast.SpecFile, error) {
	visited := make(map[string]bool)
	inProgress := make(map[string]bool)
	var result []*ast.SpecFile

	if err := r.resolveRecursive(file, basePath, visited, inProgress, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resolver) resolveRecursive(
	file *ast.SpecFile,
	basePath string,
	visited map[string]bool,
	inProgress map[string]bool,
	result *[]*ast.SpecFile,
) error {
	// 正規化されたパスを使用
	normalPath := normalizePath(basePath, file.Path)

	// サイクル検出
	if inProgress[normalPath] {
		return &errors.CycleError{Cycle: []string{normalPath}}
	}

	// 既に処理済み
	if visited[normalPath] {
		return nil
	}

	inProgress[normalPath] = true
	defer func() { delete(inProgress, normalPath) }()

	// インポートを処理
	fileDir := filepath.Dir(normalPath)
	for _, imp := range file.Imports {
		// インポートパスを解決
		impPath := normalizePath(fileDir, imp.Path)

		impFile, err := r.parser.ParseFile(impPath)
		if err != nil {
			return &errors.ImportError{
				Pos:     imp.Pos,
				Path:    imp.Path,
				Message: "failed to parse",
				Cause:   err,
			}
		}

		if err := r.resolveRecursive(impFile, filepath.Dir(impPath), visited, inProgress, result); err != nil {
			return err
		}
	}

	visited[normalPath] = true
	*result = append(*result, file)
	return nil
}

// normalizePath はパスを正規化する
func normalizePath(base, path string) string {
	if base == "." || base == "" {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(base, path))
}
