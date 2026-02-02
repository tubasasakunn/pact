// Infrastructure Config Loader - Configuration file I/O
// Source: internal/infrastructure/config/loader.go
@version("1.0")
component ConfigLoader {
	// ConfigFileName is the default config file name
	type LoaderConstants {
		configFileName: string
	}

	depends on Config
	depends on Errors

	// Configuration loading and saving
	provides LoaderAPI {
		// Create a new Loader instance
		NewLoader() -> Loader

		// Load configuration from file path
		Load(filePath: string) -> ConfigData

		// Save configuration to file path
		Save(filePath: string, cfg: ConfigData) -> error

		// Find project root by searching for .pactconfig
		FindProjectRoot(startPath: string) -> string
	}

	// Load configuration flow
	flow LoadConfig {
		fileData = self.readFile(filePath)
		if fileData == null {
			// File doesn't exist, return default config
			defaultCfg = Config.Default()
			return defaultCfg
		}
		cfg = Config.Default()
		parseResult = self.parseYAML(fileData, cfg)
		if parseResult.hasError {
			throw ConfigError
		}
		return cfg
	}

	// Save configuration flow
	flow SaveConfig {
		yamlData = self.marshalYAML(cfg)
		if yamlData == null {
			throw ConfigError
		}
		dirPath = self.dirname(filePath)
		mkdirResult = self.mkdirAll(dirPath)
		if mkdirResult.hasError {
			throw ConfigError
		}
		writeResult = self.writeFile(filePath, yamlData)
		if writeResult.hasError {
			throw ConfigError
		}
		return success
	}

	// Find project root flow
	flow FindRoot {
		absPath = self.absPath(startPath)
		if absPath == null {
			throw ConfigError
		}
		rootPath = self.searchParentDirectories(absPath, ".pactconfig")
		if rootPath == null {
			throw ConfigError
		}
		return rootPath
	}
}
