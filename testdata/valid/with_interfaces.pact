// Interfaces and implementations
component UserRepository {
	implements Repository

	provides RepositoryAPI {
		Find(id: string) -> User
		Save(user: User)
		Delete(id: string)
	}

	requires DatabaseConnection {
		Connect() -> Connection
	}
}

component OrderService {
	implements EventPublisher

	provides OrderAPI {
		CreateOrder(order: Order) -> Order
	}

	provides EventAPI {
		Publish(event: Event)
	}

	depends on OrderRepository
}
