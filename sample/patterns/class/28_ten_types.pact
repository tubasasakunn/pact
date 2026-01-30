// Pattern 28: Component with 10 different types
component MultiTypeService {
	type User {
		id: string
		name: string
		email: string
	}

	type Product {
		productId: string
		title: string
		price: float
	}

	type Order {
		orderId: string
		userId: string
		total: float
	}

	type Payment {
		paymentId: string
		amount: float
		currency: string
	}

	type Address {
		street: string
		city: string
		country: string
	}

	type Shipping {
		trackingId: string
		carrier: string
		status: string
	}

	type Review {
		reviewId: string
		rating: int
		comment: string
	}

	type Category {
		categoryId: string
		name: string
		parentId: string?
	}

	type Inventory {
		sku: string
		quantity: int
		location: string
	}

	type Notification {
		notificationId: string
		message: string
		read: bool
	}

	provides DataAPI {
		GetUser(id: string) -> User
		GetProduct(id: string) -> Product
		GetOrder(id: string) -> Order
	}
}
