// Component with flow definitions
component OrderService {
	repo: OrderRepository
	paymentService: PaymentService
	notificationService: NotificationService

	flow CreateOrder {
		start: "Receive Order Request"
		validate: "Validate Order Data"
		if valid {
			checkInventory: "Check Inventory"
			if available {
				createOrder: "Create Order Record"
				processPayment: "Process Payment"
				if paymentSuccess {
					confirmOrder: "Confirm Order"
					sendNotification: "Send Confirmation Email"
				} else {
					cancelOrder: "Cancel Order"
					notifyFailure: "Notify Payment Failure"
				}
			} else {
				notifyOutOfStock: "Notify Out of Stock"
			}
		} else {
			returnError: "Return Validation Error"
		}
		end: "Complete"
	}

	flow CancelOrder {
		start: "Receive Cancel Request"
		findOrder: "Find Order"
		if found {
			checkStatus: "Check Order Status"
			if cancellable {
				refund: "Process Refund"
				updateStatus: "Update to Cancelled"
				notify: "Send Cancellation Email"
			} else {
				rejectCancel: "Reject Cancellation"
			}
		} else {
			notFound: "Return Not Found"
		}
		end: "Complete"
	}
}
