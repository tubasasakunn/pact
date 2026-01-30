package resolver

import (
	"path/filepath"

	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
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

// Resolve はファイルリストのインポートを解決し、依存順にソートして返す
func (r *Resolver) Resolve(files []*ast.SpecFile) ([]*ast.SpecFile, error) {
	// TODO: 実装
	return nil, &errors.ImportError{Message: "not implemented"}
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
	absPath, err := filepath.Abs(filepath.Join(basePath, file.Path))
	if err != nil {
		return &errors.ImportError{Path: file.Path, Message: err.Error()}
	}

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

	// インポートを処理
	for _, imp := range file.Imports {
		impPath := filepath.Join(filepath.Dir(absPath), imp.Path)
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

	visited[absPath] = true
	*result = append(*result, file)
	return nil
}
