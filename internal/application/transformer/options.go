package transformer

// TransformOptions はクラス図変換のオプション
type TransformOptions struct {
	// FilterComponents は対象コンポーネントのフィルタリスト（空なら全て）
	FilterComponents []string
}

// SequenceOptions はシーケンス図変換のオプション
type SequenceOptions struct {
	// FlowName は変換対象のフロー名
	FlowName string
	// IncludeReturn は return イベントを含めるか
	IncludeReturn bool
}

// StateOptions は状態図変換のオプション
type StateOptions struct {
	// StatesName は変換対象の states ブロック名
	StatesName string
}

// FlowOptions はフローチャート変換のオプション
type FlowOptions struct {
	// FlowName は変換対象のフロー名
	FlowName string
	// IncludeSwimlanes はスイムレーンを含めるか
	IncludeSwimlanes bool
}
