// Component with state machine definitions
component Order {
	type OrderData {
		id: string
		total: float
	}

	states OrderStatus {
		initial Draft
		final Completed
		final Cancelled
		final Rejected
		final Refunded

		state Draft { }
		state Submitted { }
		state Confirmed { }
		state Processing { }
		state Shipped { }
		state Delivered { }
		state Returned { }
		state Completed { }
		state Cancelled { }
		state Rejected { }
		state Refunded { }

		Draft -> Submitted on submit
		Draft -> Cancelled on cancel
		Submitted -> Confirmed on confirm
		Submitted -> Rejected on reject
		Confirmed -> Processing on startProcessing
		Processing -> Shipped on ship
		Shipped -> Delivered on deliver
		Delivered -> Completed on complete
		Delivered -> Returned on returnRequest
		Returned -> Refunded on processRefund
	}
}

component Payment {
	type PaymentData {
		id: string
		amount: float
	}

	states PaymentStatus {
		initial Pending
		final Failed
		final Voided
		final Refunded

		state Pending { }
		state Authorized { }
		state Captured { }
		state PartiallyRefunded { }
		state Failed { }
		state Voided { }
		state Refunded { }

		Pending -> Authorized on authorize
		Pending -> Failed on fail
		Authorized -> Captured on capture
		Authorized -> Voided on void
		Captured -> Refunded on refund
		Captured -> PartiallyRefunded on partialRefund
		PartiallyRefunded -> Refunded on completeRefund
	}
}
