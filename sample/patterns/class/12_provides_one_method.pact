// Pattern 12: Component with provides interface (1 method)
component UserService {
	type User {
		id: string
		name: string
	}
	
	provides UserAPI {
		GetUser(id: string) -> User
	}
}
