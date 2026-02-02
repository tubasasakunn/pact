// Pattern 10: Component with depends on (1 dependency)
component Database {
	type Config {
		connectionString: string
	}
}

component Repository {
	type RepoData {
		id: string
	}
	
	depends on Database
}
