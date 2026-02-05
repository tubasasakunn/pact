// Pattern: Linear States with 4 states
// 4つの状態が直線的に遷移するパターン

component OrderProcess {
    type OrderData {
        orderId: string
        customerId: string
        totalAmount: float
    }

    states OrderState {
        initial Pending
        final Delivered

        state Pending {
            entry [createOrder]
        }

        state Confirmed {
            entry [sendConfirmation]
        }

        state Shipped {
            entry [updateTracking]
        }

        state Delivered {
            entry [completeOrder]
        }

        Pending -> Confirmed on confirm
        Confirmed -> Shipped on ship
        Shipped -> Delivered on deliver
    }
}
