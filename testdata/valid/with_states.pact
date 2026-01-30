// Component with state machine definitions
component Order {
	id: string
	status: OrderStatus
	items: []OrderItem
	total: float

	states OrderStatus {
		initial -> Draft

		Draft -> Submitted: submit
		Draft -> Cancelled: cancel

		Submitted -> Confirmed: confirm
		Submitted -> Rejected: reject
		Submitted -> Cancelled: cancel

		Confirmed -> Processing: startProcessing
		Confirmed -> Cancelled: cancel

		Processing -> Shipped: ship
		Processing -> Cancelled: cancel

		Shipped -> Delivered: deliver
		Shipped -> Returned: returnRequest

		Delivered -> Completed: complete
		Delivered -> Returned: returnRequest

		Returned -> Refunded: processRefund

		Completed -> [*]
		Cancelled -> [*]
		Rejected -> [*]
		Refunded -> [*]
	}
}

component Payment {
	id: string
	amount: float
	status: PaymentStatus

	states PaymentStatus {
		initial -> Pending

		Pending -> Authorized: authorize
		Pending -> Failed: fail

		Authorized -> Captured: capture
		Authorized -> Voided: void

		Captured -> Refunded: refund
		Captured -> PartiallyRefunded: partialRefund

		PartiallyRefunded -> Refunded: completeRefund

		Failed -> [*]
		Voided -> [*]
		Refunded -> [*]
	}
}
