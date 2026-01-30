// Pattern 10: Component with requires interface
component OrderService {
    type Order { id: string }

    requires PaymentAPI {
        ProcessPayment(amount: float) -> bool
        RefundPayment(id: string) -> bool
    }
}
