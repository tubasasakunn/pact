// 05: Simple Flow - if/else, method calls, return, throw
component PaymentProcessor {
    type Payment {
        id: string
        amount: float
        currency: string
        status: string
    }

    depends on BankAPI
    depends on FraudDetector

    flow ProcessPayment {
        fraud = FraudDetector.check(payment)
        if fraud.isSuspicious {
            throw FraudDetectedError
        }

        result = BankAPI.charge(payment.amount, payment.currency)
        if result.success {
            self.recordTransaction(payment, result)
            return result
        } else {
            throw PaymentDeclinedError
        }
    }

    flow RefundPayment {
        original = self.getPayment(paymentId)
        if original.status == "completed" {
            refund = BankAPI.refund(original.id, original.amount)
            self.updateStatus(original, "refunded")
            return refund
        } else {
            throw InvalidRefundError
        }
    }
}

component BankAPI {
    provides BankInterface {
        Charge(amount: float, currency: string) -> string
        Refund(id: string, amount: float) -> bool
    }
}

component FraudDetector {
    provides FraudAPI {
        Check(payment: string) -> string
    }
}
