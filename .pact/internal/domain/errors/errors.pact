// Domain Errors - Error types for the Pact system
// Source: internal/domain/errors/errors.go
@version("1.0")
component Errors {
	// Position in source file
	type Position {
		line: int
		column: int
		fileName: string
	}

	// ParseError represents a syntax parsing error
	type ParseError {
		pos: Position
		message: string
	}

	// SemanticError represents a semantic analysis error
	type SemanticError {
		pos: Position
		message: string
	}

	// ImportError represents an import resolution error
	type ImportError {
		pos: Position
		importPath: string
		message: string
		cause: error?
	}

	// CycleError represents a circular import detection error
	type CycleError {
		cyclePath: string[]
	}

	// TransformError represents a diagram transformation error
	type TransformError {
		sourceType: string
		targetType: string
		message: string
	}

	// ConfigError represents a configuration error
	type ConfigError {
		filePath: string
		message: string
	}

	// Error interface implementations
	provides ErrorAPI {
		// Format ParseError as string
		FormatParseError(err: ParseError) -> string

		// Format SemanticError as string
		FormatSemanticError(err: SemanticError) -> string

		// Format ImportError as string
		FormatImportError(err: ImportError) -> string

		// Format CycleError as string
		FormatCycleError(err: CycleError) -> string

		// Format TransformError as string
		FormatTransformError(err: TransformError) -> string

		// Format ConfigError as string
		FormatConfigError(err: ConfigError) -> string

		// Unwrap ImportError to get cause
		UnwrapImportError(err: ImportError) -> error
	}

	// Format import error flow
	flow FormatImportErrorFlow {
		msg = self.formatPosition(err.pos)
		msg = msg + ": import " + err.importPath + ": " + err.message
		if err.cause != null {
			msg = msg + " (caused by: " + err.cause.message + ")"
		}
		return msg
	}

	// Format cycle error flow
	flow FormatCycleErrorFlow {
		msg = "import cycle detected: ["
		for path in err.cyclePath {
			msg = msg + path + " -> "
		}
		msg = msg + "]"
		return msg
	}
}
