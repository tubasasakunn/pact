// Pattern 23: Enum type
component OrderSystem {
	enum OrderStatus {
		pending
		confirmed
		processing
		shipped
		delivered
		cancelled
		refunded
	}
	
	enum PaymentMethod {
		credit_card
		debit_card
		paypal
		bank_transfer
		cash
	}
	
	type Order {
		id: string
		status: OrderStatus
		paymentMethod: PaymentMethod
		total: float
	}
}
