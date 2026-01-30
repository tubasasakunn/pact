// Complex example with multiple features
@version("1.0")
@module("ecommerce")

// Type definitions
type Money {
	amount: float
	currency: string
}

type Email = string
type Phone = string

type ContactInfo {
	email: Email
	phone: Phone?
}

type OrderStatus = "pending" | "confirmed" | "shipped" | "delivered" | "cancelled"

// Interfaces
interface Repository[T] {
	method FindById(id: string): T?
	method FindAll(): []T
	method Save(entity: T): T
	method Delete(id: string): bool
}

interface EventHandler[E] {
	method Handle(event: E): void
}

// Domain components
@aggregate
component Customer {
	id: string
	name: string
	contact: ContactInfo
	addresses: []Address
	createdAt: timestamp

	method AddAddress(address: Address): void
	method UpdateContact(contact: ContactInfo): void
}

@aggregate
component Order {
	id: string
	customerId: string
	items: []OrderItem
	shippingAddress: Address
	billingAddress: Address
	subtotal: Money
	tax: Money
	total: Money
	status: OrderStatus
	createdAt: timestamp
	updatedAt: timestamp

	relation Customer: belongs_to
	relation OrderItem: has_many

	method AddItem(item: OrderItem): void
	method RemoveItem(itemId: string): void
	method UpdateQuantity(itemId: string, quantity: int): void
	method CalculateTotal(): Money

	states OrderLifecycle {
		initial -> Pending

		Pending -> Confirmed: confirm
		Pending -> Cancelled: cancel

		Confirmed -> Processing: process
		Confirmed -> Cancelled: cancel

		Processing -> Shipped: ship

		Shipped -> Delivered: deliver
		Shipped -> Returned: requestReturn

		Delivered -> Completed: complete
		Delivered -> Returned: requestReturn

		Returned -> Refunded: refund

		Completed -> [*]
		Cancelled -> [*]
		Refunded -> [*]
	}

	flow PlaceOrder {
		start: "Receive Order"
		validateItems: "Validate Items"
		if itemsValid {
			checkInventory: "Check Inventory"
			if inStock {
				calculatePricing: "Calculate Pricing"
				validatePayment: "Validate Payment Method"
				if paymentValid {
					authorizePayment: "Authorize Payment"
					createOrder: "Create Order"
					reserveInventory: "Reserve Inventory"
					sendConfirmation: "Send Confirmation"
					success: "Order Placed Successfully"
				} else {
					paymentError: "Return Payment Error"
				}
			} else {
				outOfStock: "Return Out of Stock Error"
			}
		} else {
			validationError: "Return Validation Error"
		}
		end: "Complete"
	}
}

@entity
component OrderItem {
	id: string
	orderId: string
	productId: string
	productName: string
	quantity: int
	unitPrice: Money
	totalPrice: Money

	relation Order: belongs_to
	relation Product: references
}

@entity
component Product {
	id: string
	sku: string
	name: string
	description: string
	price: Money
	inventory: int
	category: string
	tags: []string

	method UpdatePrice(price: Money): void
	method AdjustInventory(delta: int): void
}

// Infrastructure components
@service
component OrderService {
	implements Repository[Order]
	implements EventHandler[OrderEvent]

	orderRepo: OrderRepository
	customerRepo: CustomerRepository
	inventoryService: InventoryService
	paymentService: PaymentService
	eventPublisher: EventPublisher

	method PlaceOrder(request: PlaceOrderRequest): Order
	method CancelOrder(orderId: string): void
	method GetOrderStatus(orderId: string): OrderStatus
	method Handle(event: OrderEvent): void
}

@repository
component OrderRepository {
	implements Repository[Order]

	db: Database
	cache: CacheService

	method FindById(id: string): Order?
	method FindAll(): []Order
	method FindByCustomer(customerId: string): []Order
	method Save(order: Order): Order
	method Delete(id: string): bool
}
