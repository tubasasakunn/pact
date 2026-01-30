// Interfaces and implementations
interface Repository {
	method Find(id: string): Entity
	method Save(entity: Entity): void
	method Delete(id: string): void
}

interface EventPublisher {
	method Publish(event: Event): void
}

component UserRepository {
	implements Repository

	db: Database

	method Find(id: string): User
	method Save(user: User): void
	method Delete(id: string): void
}

component OrderService {
	implements EventPublisher

	repo: OrderRepository
	publisher: EventPublisher

	method CreateOrder(order: Order): Order
	method Publish(event: Event): void
}
