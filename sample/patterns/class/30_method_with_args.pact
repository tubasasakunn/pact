// Pattern 30: Methods with multiple arguments
component DataProcessor {
	type ProcessResult {
		success: bool
		message: string
		processedCount: int
	}

	type FilterCriteria {
		field: string
		operator: string
		value: string
	}

	type SortOptions {
		field: string
		ascending: bool
	}

	type PageInfo {
		page: int
		pageSize: int
		totalPages: int
		totalItems: int
	}

	provides ProcessingAPI {
		// Single argument
		ProcessSingle(data: string) -> ProcessResult

		// Two arguments
		ProcessPair(input: string, output: string) -> ProcessResult

		// Three arguments
		ProcessTriple(source: string, target: string, format: string) -> ProcessResult

		// Four arguments
		TransformData(input: string, output: string, format: string, validate: bool) -> ProcessResult

		// Five arguments
		ComplexProcess(source: string, dest: string, filter: FilterCriteria, sort: SortOptions, limit: int) -> ProcessResult

		// Six arguments with mixed types
		BatchOperation(items: string[], batchSize: int, runParallel: bool, retryCount: int, timeout: int, callback: string) -> ProcessResult

		// Many arguments for pagination
		SearchWithPagination(query: string, filters: FilterCriteria[], sort: SortOptions, page: int, pageSize: int, includeDeleted: bool) -> PageInfo
	}
}
