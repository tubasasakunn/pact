// Components with relations
component User {
	id: string
	name: string
}

component Order {
	id: string
	userId: string
	total: float

	relation User: uses
}

component OrderItem {
	id: string
	orderId: string
	productId: string
	quantity: int

	relation Order: belongs_to
	relation Product: references
}

component Product {
	id: string
	name: string
	price: float
}
