// Component with annotations
@version("2.0")
@author("development-team")
@deprecated("Use UserServiceV2 instead")
component UserService {
	@inject
	private db: Database

	@cache(ttl: 300)
	private cache: CacheService

	@log(level: "debug")
	@metric(name: "user_service_requests")
	public method GetUser(id: string): User

	@async
	@retry(maxAttempts: 3, delay: 1000)
	public method CreateUser(user: User): User

	@transactional
	@audit(action: "delete_user")
	public method DeleteUser(id: string): void
}

@api(version: "v1", basePath: "/api/v1")
interface UserAPI {
	@get("/users/{id}")
	@auth(required: true)
	method GetUser(id: string): User

	@post("/users")
	@validate
	method CreateUser(user: CreateUserRequest): User

	@delete("/users/{id}")
	@auth(role: "admin")
	method DeleteUser(id: string): void
}
