package export

import "io"

// ExportFormat はエクスポート形式
type ExportFormat string

const (
	FormatSVG ExportFormat = "svg"
	FormatPNG ExportFormat = "png"
	FormatPDF ExportFormat = "pdf"
)

// Exporter はダイアグラムをエクスポートするインターフェース
type Exporter interface {
	// Export はダイアグラムを指定形式でエクスポートする
	Export(svgData []byte, format ExportFormat, w io.Writer) error
	// SupportedFormats はサポートする形式のリストを返す
	SupportedFormats() []ExportFormat
}

// SVGExporter はSVGエクスポーター（デフォルト）
type SVGExporter struct{}

// NewSVGExporter は新しいSVGエクスポーターを作成する
func NewSVGExporter() *SVGExporter {
	return &SVGExporter{}
}

// Export はSVGをそのまま出力する
func (e *SVGExporter) Export(svgData []byte, format ExportFormat, w io.Writer) error {
	if format != FormatSVG {
		return &ExportError{Format: string(format), Message: "format not supported (only SVG is available; PNG/PDF requires external converter)"}
	}
	_, err := w.Write(svgData)
	return err
}

// SupportedFormats はサポートする形式のリストを返す
func (e *SVGExporter) SupportedFormats() []ExportFormat {
	return []ExportFormat{FormatSVG}
}

// ExportError はエクスポートエラー
type ExportError struct {
	Format  string
	Message string
}

func (e *ExportError) Error() string {
	return "export " + e.Format + ": " + e.Message
}
