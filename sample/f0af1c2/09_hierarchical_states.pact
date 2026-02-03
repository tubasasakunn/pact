// 09: Hierarchical (Nested) States
component ATMController {
    type Transaction {
        id: string
        amount: float
        kind: string
    }

    states ATMOperation {
        initial Idle
        final OutOfService

        state Idle {
            entry [displayWelcome]
        }

        state Active {
            initial Authenticating

            state Authenticating {
                entry [requestCard]
            }

            state Processing {
                entry [showMenu]
            }

            Authenticating -> Processing on pinVerified
        }

        Idle -> Active on cardInserted
        Active -> Idle on ejectCard
        Active -> Idle after 30s
        Idle -> OutOfService on maintenance
    }
}
