// Pattern 22: Mixed relationships (depends, extends, implements)
component BaseService {
	type BaseConfig {
		name: string
	}
}

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

component AdvancedService {
	type ServiceData {
		id: string
		status: string
	}
	
	extends BaseService
	implements Cacheable
	implements Loggable
	depends on Logger
	depends on Cache
}
