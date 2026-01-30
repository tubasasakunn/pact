// Component with annotations
@version("2.0")
@author("development-team")
component UserService {
	type User {
		id: string
		name: string
	}

	depends on Database
	depends on CacheService

	@cache(ttl: "300")
	provides UserAPI {
		@log(level: "debug")
		GetUser(id: string) -> User

		@async
		CreateUser(user: User) -> User

		@transactional
		DeleteUser(id: string)
	}
}
