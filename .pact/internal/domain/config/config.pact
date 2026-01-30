// Domain Config - Project configuration model
// Source: internal/domain/config/config.go
@version("1.0")
component Config {
	// Config represents the project configuration
	type ConfigData {
		sourceRoot: string
		pactRoot: string
		outputDir: string
		language: string
		diagramList: string[]
		excludePatterns: string[]
	}

	// Configuration management interface
	provides ConfigAPI {
		// Create default configuration
		Default() -> ConfigData

		// Check if a diagram type is enabled
		DiagramEnabled(diagramKind: string) -> bool

		// Check if a path should be excluded
		IsExcluded(filePath: string) -> bool
	}

	// Check diagram enabled flow
	flow CheckDiagramEnabled {
		for diagramKind in self.diagramList {
			if diagramKind == "all" {
				return true
			}
			if diagramKind == targetType {
				return true
			}
		}
		return false
	}

	// Check exclusion flow
	flow CheckExclusion {
		for pattern in self.excludePatterns {
			matched = self.matchPattern(pattern, filePath)
			if matched {
				return true
			}
			baseMatched = self.matchPattern(pattern, self.baseName(filePath))
			if baseMatched {
				return true
			}
		}
		return false
	}
}
