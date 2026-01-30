// Infrastructure Import Resolver - Import and dependency resolution
// Source: internal/infrastructure/resolver/import.go
@version("1.0")
component ImportResolver {
	// Parser interface for file parsing
	type ParserInterface {
		parseFile: func
	}

	depends on Parser
	depends on Errors

	// Import resolution interface
	provides ResolverAPI {
		// Create a new Resolver with parser
		NewResolver(parser: ParserInterface) -> Resolver

		// Create a new ImportResolver
		NewImportResolver() -> ImportResolver

		// Resolve imports from a file path, returns dependency-ordered paths
		Resolve(filePath: string) -> string[]

		// Resolve imports for a list of spec files
		ResolveFiles(files: SpecFile[]) -> SpecFile[]

		// Resolve imports for a single file with base path
		ResolveFile(specFile: SpecFile, basePath: string) -> SpecFile[]

		// Check if error is a cycle error
		IsCycleError(err: error) -> bool
	}

	// Resolve single file imports flow
	flow ResolveImports {
		absPath = self.absPath(filePath)
		if absPath == null {
			throw ImportError
		}
		visited = self.newVisitedSet()
		inProgress = self.newInProgressSet()
		resultList = self.newResultList()
		self.resolveFileRecursive(absPath, visited, inProgress, resultList)
		return resultList
	}

	// Recursive file resolution flow
	flow ResolveFileRecursive {
		hasCycle = self.checkCycle(inProgress, absPath)
		if hasCycle {
			throw CycleError
		}
		isVisited = self.checkVisited(visited, absPath)
		if isVisited {
			return success
		}
		self.markInProgress(inProgress, absPath)
		content = self.readFile(absPath)
		if content == null {
			throw ImportError
		}
		spec = self.parseContent(content)
		self.processImports(spec, visited, inProgress, resultList)
		self.markVisited(visited, absPath)
		self.addToResult(resultList, absPath)
		return success
	}

	states ImportResolution {
		initial Pending
		final Resolved
		final Failed

		state Pending { }
		state InProgress { }
		state Resolved { }
		state Failed { }

		Pending -> InProgress on startResolve
		InProgress -> Resolved on resolveSuccess
		InProgress -> Failed on cycleDetected
		InProgress -> Failed on fileNotFound
		InProgress -> Failed on parseError
	}
}
