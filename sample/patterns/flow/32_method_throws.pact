// Pattern 32: Method throws declarations
component PaymentService {
    depends on PaymentGateway : Gateway as gateway
    depends on AccountService : Account as account

    type PaymentError {
        code: string
        message: string
    }

    provides PaymentAPI {
        processPayment(amount: float, cardId: string) -> Receipt throws PaymentError, NetworkError
        refund(transactionId: string) -> bool throws RefundError
        validateCard(cardNumber: string) -> bool throws ValidationError, ExpiredCardError
        async transferFunds(from: string, to: string, amount: float) -> Transfer throws InsufficientFundsError, AccountLockedError
    }

    requires Logging {
        log(level: string, message: string)
        logError(error: PaymentError) throws LoggingError
    }

    flow ProcessPaymentWithErrors {
        validated = self.validateCard(cardNumber)
        if validated {
            balance = account.getBalance(accountId)
            if balanceInsufficient {
                throw InsufficientFundsError
            }
            result = gateway.charge(amount, cardId)
            return result
        } else {
            throw ValidationError
        }
    }
}
