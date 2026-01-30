// Real-world example: Order Service
@version("1.0")
component OrderService {
	type Money {
		amount: float
		currency: string
	}

	type Address {
		line1: string
		city: string
		postalCode: string
		country: string
	}

	type OrderItem {
		productId: string
		quantity: int
		unitPrice: Money
	}

	depends on OrderRepository
	depends on PaymentGateway
	depends on InventoryService
	depends on ShippingService
	depends on NotificationService

	provides OrderAPI {
		CreateOrder(request: CreateOrderRequest) -> Order
		SubmitOrder(orderId: string) -> Order
		ProcessPayment(orderId: string) -> PaymentResult
		CancelOrder(orderId: string, reason: string)
		ShipOrder(orderId: string) -> ShipmentResult
		GetOrderStatus(orderId: string) -> OrderStatus
	}

	states OrderStatus {
		initial Draft
		final Completed
		final Cancelled
		final Refunded

		state Draft { }
		state Submitted { }
		state PaymentPending { }
		state PaymentAuthorized { }
		state Processing { }
		state ReadyToShip { }
		state Shipped { }
		state Delivered { }
		state Completed { }
		state Cancelled { }
		state Refunded { }

		Draft -> Submitted on submit
		Submitted -> PaymentPending on awaitPayment
		PaymentPending -> PaymentAuthorized on authorizePayment
		PaymentPending -> Cancelled on cancel
		PaymentAuthorized -> Processing on startProcessing
		Processing -> ReadyToShip on packComplete
		ReadyToShip -> Shipped on ship
		Shipped -> Delivered on deliver
		Delivered -> Completed on complete
	}

	flow SubmitOrder {
		valid = self.validateOrder(order)
		if valid {
			available = InventoryService.checkAvailability(order.items)
			if available {
				reserved = InventoryService.reserve(order.id, order.items)
				if reserved {
					self.calculateTotals(order)
					authorized = PaymentGateway.authorize(order.total, order.paymentInfo)
					if authorized {
						saved = OrderRepository.save(order)
						self.publishEvent(saved)
						NotificationService.sendConfirmation(saved)
						return saved
					} else {
						InventoryService.release(order.id)
						throw PaymentError
					}
				} else {
					throw InventoryError
				}
			} else {
				throw OutOfStockError
			}
		} else {
			throw ValidationError
		}
	}
}
