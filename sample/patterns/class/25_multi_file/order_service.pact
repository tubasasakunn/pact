// Pattern 25c: Order service referencing shared and user
import "./shared.pact"
import "./user_service.pact"

component OrderService {
	type Order {
		id: string
		userId: string
		total: float
		status: Status
	}
	
	depends on SharedTypes
	depends on UserService
	
	provides OrderAPI {
		GetOrder(id: string) -> Order
		CreateOrder(order: Order) -> Order
		GetUserOrders(userId: string) -> Order[]
	}
}
