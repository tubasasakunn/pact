// Pattern 25b: User service referencing shared
import "./shared.pact"

component UserService {
	type User {
		id: string
		name: string
		email: string
		status: Status
	}
	
	depends on SharedTypes
	
	provides UserAPI {
		GetUser(id: string) -> User
		CreateUser(user: User) -> User
		UpdateUser(id: string, user: User) -> User
	}
}
