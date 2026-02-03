package watcher

// WatchMode はファイル監視モードの設定
type WatchMode struct {
	Enabled    bool
	Paths      []string // 監視対象パス
	Extensions []string // 監視対象の拡張子（デフォルト: [".pact"]）
}

// WatchEvent はファイル変更イベント
type WatchEvent struct {
	Path string
	Type WatchEventType
}

// WatchEventType はイベントの種類
type WatchEventType int

const (
	WatchEventModified WatchEventType = iota
	WatchEventCreated
	WatchEventDeleted
)

// Watcher はファイル監視のインターフェース
type Watcher interface {
	// Watch はファイル監視を開始する
	Watch(config WatchMode) (<-chan WatchEvent, error)
	// Stop はファイル監視を停止する
	Stop() error
}

// NewWatchMode はデフォルトのWatchMode設定を返す
func NewWatchMode(paths ...string) WatchMode {
	return WatchMode{
		Enabled:    true,
		Paths:      paths,
		Extensions: []string{".pact"},
	}
}
