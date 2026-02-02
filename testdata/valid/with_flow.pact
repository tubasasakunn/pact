// Component with flow definitions
component OrderService {
	depends on OrderRepository
	depends on PaymentService
	depends on NotificationService

	flow CreateOrder {
		valid = self.validateOrder(request)
		if valid {
			available = self.checkInventory(items)
			if available {
				order = OrderRepository.create(request)
				paymentSuccess = PaymentService.process(order)
				if paymentSuccess {
					self.confirmOrder(order)
					NotificationService.sendConfirmation(order)
				} else {
					self.cancelOrder(order)
					NotificationService.notifyFailure(order)
				}
			} else {
				NotificationService.notifyOutOfStock(items)
			}
		} else {
			throw ValidationError
		}
		return order
	}

	flow CancelOrder {
		order = OrderRepository.find(orderId)
		if order {
			cancellable = self.checkCancellable(order)
			if cancellable {
				PaymentService.refund(order)
				self.updateStatus(order)
				NotificationService.sendCancellation(order)
			} else {
				throw CancellationRejected
			}
		} else {
			throw OrderNotFound
		}
	}
}
