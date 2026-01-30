// Complex example with multiple features
@version("1.0")
component Order {
	type Money {
		amount: float
		currency: string
	}

	type OrderData {
		id: string
		customerId: string
		total: Money
	}

	depends on Customer
	depends on OrderItem

	provides OrderAPI {
		AddItem(item: OrderItem)
		RemoveItem(itemId: string)
		CalculateTotal() -> Money
	}

	states OrderLifecycle {
		initial Pending
		final Completed
		final Cancelled
		final Refunded

		state Pending { }
		state Confirmed { }
		state Processing { }
		state Shipped { }
		state Delivered { }
		state Completed { }
		state Cancelled { }
		state Refunded { }

		Pending -> Confirmed on confirm
		Pending -> Cancelled on cancel
		Confirmed -> Processing on process
		Processing -> Shipped on ship
		Shipped -> Delivered on deliver
		Delivered -> Completed on complete
	}

	flow PlaceOrder {
		valid = self.validateItems(items)
		if valid {
			inStock = self.checkInventory(items)
			if inStock {
				total = self.calculatePricing(items)
				paymentValid = self.validatePayment(payment)
				if paymentValid {
					self.authorizePayment(payment)
					order = self.createOrder(items)
					self.reserveInventory(items)
					self.sendConfirmation(order)
					return order
				} else {
					throw PaymentError
				}
			} else {
				throw OutOfStockError
			}
		} else {
			throw ValidationError
		}
	}
}

component Customer {
	type CustomerData {
		id: string
		name: string
	}
}

component OrderItem {
	type ItemData {
		id: string
		productId: string
		quantity: int
	}

	depends on Product
}

component Product {
	type ProductData {
		id: string
		name: string
		price: float
	}
}
