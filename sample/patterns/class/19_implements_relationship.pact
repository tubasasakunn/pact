// Pattern 19: Component with implements relationship
component UserRepository {
	type User {
		id: string
		name: string
	}
	
	implements Repository
	
	provides RepositoryAPI {
		Find(id: string) -> User
		Save(user: User) -> User
		Delete(id: string)
	}
}
