// 03: Provides/Requires Interfaces with async, throws
component ShoppingService {
    type CartItem {
        productId: string
        quantity: int
        price: float
    }

    type Cart {
        userId: string
        items: CartItem[]
        total: float
    }

    provides CartAPI {
        AddItem(userId: string, item: CartItem) -> Cart
        RemoveItem(userId: string, productId: string) -> Cart
        GetCart(userId: string) -> Cart
        Checkout(userId: string) -> string throws EmptyCartError
        async CalculateShipping(cart: Cart) -> float
    }

    requires InventoryService {
        CheckStock(productId: string) -> int
        ReserveStock(productId: string, qty: int) -> bool throws OutOfStockError
    }

    requires PaymentService {
        async ProcessPayment(amount: float, method: string) -> string throws PaymentError
        RefundPayment(transactionId: string) -> bool
    }
}
