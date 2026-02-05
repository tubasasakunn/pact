// Pattern: Binary Choice
// 分岐を含む状態遷移パターン

component PaymentProcess {
    type PaymentData {
        paymentId: string
        amount: float
        method: string
    }

    states PaymentState {
        initial Processing
        final Completed
        final Failed

        state Processing {
            entry [initiatePayment]
        }

        state Completed {
            entry [recordSuccess, notifyCustomer]
        }

        state Failed {
            entry [recordFailure, notifySupport]
        }

        Processing -> Completed on success
        Processing -> Failed on error
    }
}
