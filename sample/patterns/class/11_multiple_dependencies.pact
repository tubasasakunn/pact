// Pattern 11: Component with multiple depends on
component Logger {
	type LogEntry {
		message: string
	}
}

component Cache {
	type CacheEntry {
		key: string
		value: string
	}
}

component Database {
	type Config {
		connectionString: string
	}
}

component Service {
	type ServiceData {
		id: string
	}
	
	depends on Logger
	depends on Cache
	depends on Database
}
