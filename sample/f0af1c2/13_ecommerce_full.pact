// 13: Full E-commerce System - all features combined
@version("3.0.0")
@author("ecommerce-team")
@description("Complete e-commerce platform")
component ProductCatalog {
    @entity
    type Product {
        +id: string
        +name: string
        +description: string
        +price: float
        +categories: string[]
        +inStock: bool
        -costPrice: float
    }

    enum ProductStatus {
        Draft
        Active
        Discontinued
        OutOfStock
    }

    provides CatalogAPI {
        SearchProducts(query: string) -> Product[]
        GetProduct(id: string) -> Product
        ListByCategory(category: string) -> Product[]
    }
}

component ShoppingCart {
    type CartItem {
        productId: string
        quantity: int
        unitPrice: float
    }

    type Cart {
        userId: string
        items: CartItem[]
        totalAmount: float
    }

    depends on ProductCatalog

    provides CartAPI {
        AddItem(userId: string, productId: string, qty: int) -> Cart throws OutOfStockError
        RemoveItem(userId: string, productId: string) -> Cart
        GetCart(userId: string) -> Cart
        ClearCart(userId: string) -> bool
    }

    flow AddToCart {
        product = ProductCatalog.GetProduct(productId)
        if product.inStock {
            cart = self.getOrCreateCart(userId)
            self.addItemToCart(cart, product, quantity)
            total = self.recalculateTotal(cart)
            return cart
        } else {
            throw OutOfStockError
        }
    }
}

component OrderManager {
    type Order {
        +id: string
        +userId: string
        +items: string[]
        +total: float
        +status: string
        +createdAt: string
    }

    depends on ShoppingCart
    depends on PaymentGateway
    depends on Warehouse

    provides OrderAPI {
        PlaceOrder(userId: string) -> Order throws PaymentError
        CancelOrder(orderId: string) -> bool
        GetOrderStatus(orderId: string) -> string
    }

    flow PlaceOrder {
        cart = ShoppingCart.GetCart(userId)
        if cart.totalAmount > 0 {
            payment = PaymentGateway.charge(userId, cart.totalAmount)
            if payment.success {
                order = self.createOrder(userId, cart)
                Warehouse.allocateStock(order)
                ShoppingCart.ClearCart(userId)
                return order
            } else {
                throw PaymentError
            }
        } else {
            throw EmptyCartError
        }
    }

    states OrderLifecycle {
        initial Pending
        final Delivered
        final Cancelled
        final Refunded

        state Pending {
            entry [logOrderCreated]
        }

        state Confirmed {
            entry [sendConfirmationEmail]
        }

        state Processing {
            entry [notifyWarehouse]
        }

        state Shipped {
            entry [sendTrackingInfo]
        }

        state Delivered {
            entry [requestFeedback]
        }

        state Cancelled {
            entry [processRefund, restoreStock]
        }

        state Refunded {
            entry [confirmRefund]
        }

        Pending -> Confirmed on paymentVerified
        Confirmed -> Processing on startFulfillment
        Processing -> Shipped on dispatched
        Shipped -> Delivered on delivered
        Pending -> Cancelled on cancel
        Confirmed -> Cancelled on cancel
        Delivered -> Refunded on refundRequested when withinReturnPeriod
    }
}

component PaymentGateway {
    type PaymentResult {
        success: bool
        transactionId: string
        errorMessage: string?
    }

    provides PaymentAPI {
        Charge(userId: string, amount: float) -> PaymentResult
        Refund(transactionId: string) -> bool
    }
}

component Warehouse {
    type StockRecord {
        productId: string
        available: int
        reserved: int
    }

    provides WarehouseAPI {
        AllocateStock(order: string) -> bool
        ReleaseStock(order: string) -> bool
        GetStockLevel(productId: string) -> int
    }
}
