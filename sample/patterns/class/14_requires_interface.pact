// Pattern 14: Component with requires interface
component OrderProcessor {
	type Order {
		id: string
		total: float
	}
	
	requires PaymentGateway {
		ProcessPayment(amount: float) -> bool
		RefundPayment(transactionId: string) -> bool
	}
}
