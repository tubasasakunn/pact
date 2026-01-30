// Pattern 24: Complex component with all features
@version("3.0.0")
@author("platform-team")
@description("Complete e-commerce order management system")
component OrderManagement {
	type Money {
		amount: float
		currency: string
	}
	
	type Address {
		street: string
		city: string
		state: string
		zipCode: string
		country: string
	}
	
	type OrderItem {
		productId: string
		productName: string
		quantity: int
		unitPrice: Money
		discount: float?
	}
	
	type Order {
		id: string
		customerId: string
		items: OrderItem[]
		shippingAddress: Address
		billingAddress: Address?
		subtotal: Money
		tax: Money
		total: Money
		notes: string?
	}
	
	enum OrderStatus {
		draft
		pending
		confirmed
		processing
		shipped
		delivered
		cancelled
		refunded
	}
	
	extends BaseEntity
	implements Auditable
	implements Versionable
	
	depends on CustomerService
	depends on InventoryService
	depends on PaymentService
	depends on NotificationService
	
	contains OrderItemProcessor
	aggregates ShippingHandler
	
	provides OrderAPI {
		@cache(ttl: "60")
		GetOrder(id: string) -> Order
		
		CreateOrder(order: Order) -> Order
		
		@async
		UpdateOrder(id: string, order: Order) -> Order
		
		@transactional
		CancelOrder(id: string)
		
		ListOrders(customerId: string) -> Order[]
		
		GetOrderStatus(id: string) -> OrderStatus
	}
	
	provides AdminAPI {
		@admin_only
		ForceComplete(id: string)
		
		@admin_only
		RefundOrder(id: string, reason: string)
	}
	
	requires PaymentGateway {
		ProcessPayment(amount: Money) -> string
		RefundPayment(transactionId: string) -> bool
	}
	
	requires ShippingProvider {
		CreateShipment(order: Order) -> string
		TrackShipment(trackingId: string) -> string
	}
}

component CustomerService {
	type Customer {
		id: string
		name: string
	}
}

component InventoryService {
	type Stock {
		productId: string
		quantity: int
	}
}

component PaymentService {
	type Payment {
		id: string
		amount: float
	}
}

component NotificationService {
	type Notification {
		message: string
	}
}

component OrderItemProcessor {
	type ProcessedItem {
		itemId: string
		status: string
	}
}

component ShippingHandler {
	type Shipment {
		trackingId: string
	}
}

component BaseEntity {
	type BaseData {
		createdAt: string
		updatedAt: string
	}
}
