package theme

// Theme はダイアグラムのテーマ設定
type Theme struct {
	Name            string
	BackgroundColor string
	NodeFill        string
	NodeStroke      string
	NodeTextColor   string
	EdgeColor       string
	LabelColor      string
	NoteFill        string
	NoteStroke      string
	FontFamily      string
	FontSize        int
}

// DefaultTheme はデフォルトテーマ
func DefaultTheme() *Theme {
	return &Theme{
		Name:            "default",
		BackgroundColor: "#ffffff",
		NodeFill:        "#ffffff",
		NodeStroke:      "#000000",
		NodeTextColor:   "#000000",
		EdgeColor:       "#000000",
		LabelColor:      "#333333",
		NoteFill:        "#ffffcc",
		NoteStroke:      "#cccc00",
		FontFamily:      "monospace",
		FontSize:        12,
	}
}

// DarkTheme はダークテーマ
func DarkTheme() *Theme {
	return &Theme{
		Name:            "dark",
		BackgroundColor: "#1e1e1e",
		NodeFill:        "#2d2d2d",
		NodeStroke:      "#888888",
		NodeTextColor:   "#d4d4d4",
		EdgeColor:       "#888888",
		LabelColor:      "#cccccc",
		NoteFill:        "#3d3d1e",
		NoteStroke:      "#666633",
		FontFamily:      "monospace",
		FontSize:        12,
	}
}

// BlueprintTheme はブループリントテーマ
func BlueprintTheme() *Theme {
	return &Theme{
		Name:            "blueprint",
		BackgroundColor: "#1a237e",
		NodeFill:        "#283593",
		NodeStroke:      "#5c6bc0",
		NodeTextColor:   "#e8eaf6",
		EdgeColor:       "#7986cb",
		LabelColor:      "#c5cae9",
		NoteFill:        "#1a237e",
		NoteStroke:      "#5c6bc0",
		FontFamily:      "monospace",
		FontSize:        12,
	}
}

// GetTheme は名前からテーマを返す
func GetTheme(name string) *Theme {
	switch name {
	case "dark":
		return DarkTheme()
	case "blueprint":
		return BlueprintTheme()
	default:
		return DefaultTheme()
	}
}
