package renderer

// LSPCapabilities はLanguage Server Protocolのケイパビリティ
type LSPCapabilities struct {
	// Completion はオートコンプリート対応
	Completion bool
	// Diagnostics はリアルタイム診断対応
	Diagnostics bool
	// Hover はホバー情報対応
	Hover bool
	// GoToDefinition は定義ジャンプ対応
	GoToDefinition bool
	// References は参照検索対応
	References bool
	// Formatting はフォーマット対応
	Formatting bool
}

// LSPServer はLanguage Serverのインターフェース
type LSPServer interface {
	// Initialize はLSPサーバーを初期化する
	Initialize() (*LSPCapabilities, error)
	// Shutdown はLSPサーバーを停止する
	Shutdown() error
}

// DefaultLSPCapabilities はデフォルトのLSPケイパビリティを返す
func DefaultLSPCapabilities() *LSPCapabilities {
	return &LSPCapabilities{
		Completion:     true,
		Diagnostics:    true,
		Hover:          true,
		GoToDefinition: true,
		References:     true,
		Formatting:     true,
	}
}
