// Pattern 12: Enum type
component OrderService {
    enum OrderStatus {
        Pending
        Processing
        Shipped
        Delivered
        Cancelled
    }

    type Order {
        id: string
        status: OrderStatus
    }
}
