// Pattern 08: Multiple dependencies
component OrderService {
    type Order { id: string }
    depends on UserService
    depends on ProductService
    depends on PaymentService
    depends on NotificationService
}

component UserService { }
component ProductService { }
component PaymentService { }
component NotificationService { }
