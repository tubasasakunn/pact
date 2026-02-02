// Pattern 9: Conditional call (if)
component PaymentProcessor {
    depends on CardService
    depends on BankService

    flow ProcessPayment {
        if paymentType == "card" {
            result = CardService.charge(amount)
        } else {
            result = BankService.transfer(amount)
        }
        return result
    }
}

component CardService { }
component BankService { }
