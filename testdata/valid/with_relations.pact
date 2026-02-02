// Components with relations
component User {
	type UserData {
		id: string
		name: string
	}
}

component Order {
	type OrderData {
		id: string
		userId: string
		total: float
	}

	depends on User
}

component OrderItem {
	type OrderItemData {
		id: string
		orderId: string
		productId: string
		quantity: int
	}

	depends on Order
	depends on Product
}

component Product {
	type ProductData {
		id: string
		name: string
		price: float
	}
}
