# Pact DSL Specification

Pact is a domain-specific language for defining software architecture and generating diagrams.

## Overview

Pact supports four diagram types:
- **Class Diagrams**: Component structure, types, and relationships
- **Sequence Diagrams**: Message flows between components
- **Flow Diagrams**: Flowcharts with control flow
- **State Diagrams**: State machines and transitions

## File Structure

```pact
// Imports
import "path/to/module"
import "path/to/other" as alias

// Component definition
@annotation("value")
component ComponentName {
    // Types
    type TypeName { ... }
    enum EnumName { ... }

    // Relations
    depends on OtherComponent
    extends BaseComponent
    implements Interface
    contains ChildComponent
    aggregates Reference

    // Interfaces
    provides API { ... }
    requires ExternalAPI { ... }

    // Flows
    flow ProcessName { ... }

    // State machines
    states StateMachineName { ... }
}
```

## Token Types

### Keywords

| Category | Keywords |
|----------|----------|
| Structure | `component`, `import`, `as` |
| Types | `type`, `enum` |
| Relations | `depends`, `on`, `extends`, `implements`, `contains`, `aggregates` |
| Interfaces | `provides`, `requires`, `async`, `throws` |
| Flows | `flow`, `return`, `throw`, `if`, `else`, `for`, `in`, `while`, `await` |
| States | `states`, `state`, `parallel`, `region`, `initial`, `final`, `entry`, `exit`, `when`, `after`, `do` |
| Literals | `true`, `false`, `null` |

### Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+`, `-`, `*`, `/`, `%` |
| Comparison | `==`, `!=`, `<`, `>`, `<=`, `>=` |
| Logical | `&&`, `\|\|`, `!` |
| Other | `=`, `->`, `.`, `?`, `??`, `@` |

### Visibility Modifiers

| Symbol | Visibility |
|--------|------------|
| `+` | Public |
| `-` | Private |
| `#` | Protected |
| `~` | Package |

## Grammar

### Type Declarations

```pact
// Struct type
type User {
    +id: string           // public
    -password: string     // private
    #email: string        // protected
    ~internal: int        // package
    roles: string[]       // array type
    manager: User?        // nullable type
}

// Enum type
enum Status {
    Active
    Inactive
    Pending
}
```

### Type Expressions

```
TypeExpr = TypeName ('?' | '[]')?
```

- `Type` - Base type
- `Type?` - Nullable type
- `Type[]` - Array type

### Relations

```pact
// Dependency with type qualifier
depends on Database: database
depends on API: external as apiClient

// Inheritance
extends BaseComponent

// Interface implementation
implements Serializable

// Composition (strong ownership)
contains Engine

// Aggregation (weak reference)
aggregates Logger
```

### Interfaces

```pact
provides UserService {
    // Sync method
    getUser(id: string) -> User

    // Async method
    async fetchUsers() -> User[]

    // Method with throws
    updateUser(user: User) -> User throws ValidationError, NotFoundError
}

requires ExternalAPI {
    fetch(url: string) -> Response
}
```

### Flows

```pact
flow ProcessOrder {
    // Assignment
    order = orderService.getOrder(orderId)

    // Conditional
    if order.isValid() {
        result = paymentService.charge(order)
    } else {
        throw InvalidOrderError
    }

    // Loop
    for item in order.items {
        inventory.reserve(item)
    }

    // While loop
    while !order.isComplete() {
        order.processNext()
    }

    // Async call
    await notificationService.send(order)

    // Return
    return order
}
```

### Expressions

```pact
// Literals
42                          // int
3.14                        // float
"hello"                     // string
true                        // boolean
null                        // null

// Variable
user

// Field access
user.name

// Method call
user.getName()
service.process(arg1, arg2)

// Binary operations
a + b
x > 10 && y < 20

// Unary operations
!isValid
-amount

// Ternary
condition ? thenExpr : elseExpr

// Nullish coalescing
value ?? defaultValue
value ?? throw NullError
```

### State Machines

```pact
states OrderState {
    initial Pending
    final Completed
    final Cancelled

    // Simple state
    state Processing {
        entry [startProcessing]
        exit [cleanup]
    }

    // Hierarchical state
    state Active {
        initial SubState1

        state SubState1 { }
        state SubState2 { }

        SubState1 -> SubState2 on advance
    }

    // Transitions
    Pending -> Processing on submit
    Processing -> Completed when order.isPaid()
    Processing -> Cancelled after 24h

    // Transition with guard and actions
    Processing -> Shipped on ship when inventory.available() do [updateInventory, notify]

    // Parallel state
    parallel ActiveState {
        region BillingRegion {
            initial Unbilled
            state Unbilled { }
            state Billed { }
            Unbilled -> Billed on charge
        }
        region ShippingRegion {
            initial Preparing
            state Preparing { }
            state Shipped { }
            Preparing -> Shipped on ship
        }
    }
}
```

### Annotations

```pact
@deprecated
@note("This is a description")
@custom(key: "value", other: "data")
component Example { }
```

## AST Structure

### Root

```
SpecFile
├── Imports[]
├── Component (single, legacy)
├── Components[] (multiple)
├── Interfaces[]
├── Types[]
└── Annotations[]
```

### Component

```
ComponentDecl
├── Name: string
├── Annotations[]
└── Body
    ├── Types[]
    ├── Relations[]
    ├── Provides[]
    ├── Requires[]
    ├── Flows[]
    └── States[]
```

### Expression Hierarchy

```
Expr (interface)
├── LiteralExpr (int, float, string, bool, null)
├── VariableExpr
├── FieldExpr
├── CallExpr
├── BinaryExpr
├── UnaryExpr
├── TernaryExpr
└── NullishExpr
```

### Step Hierarchy

```
Step (interface)
├── AssignStep
├── CallStep
├── ReturnStep
├── ThrowStep
├── IfStep
├── ForStep
└── WhileStep
```

## Operator Precedence

| Level | Operators | Associativity |
|-------|-----------|---------------|
| 1 | `?:`, `??` | Right |
| 2 | `\|\|` | Left |
| 3 | `&&` | Left |
| 4 | `==`, `!=` | Left |
| 5 | `<`, `>`, `<=`, `>=` | Left |
| 6 | `+`, `-` | Left |
| 7 | `*`, `/`, `%` | Left |
| 8 | `!`, `-` (unary) | Right |
| 9 | `.` (member) | Left |

## Duration Literals

```pact
500ms   // milliseconds
30s     // seconds
5m      // minutes
24h     // hours
7d      // days
```

## Comments

```pact
// Single line comment

/*
  Multi-line
  comment
*/
```

## String Escapes

| Escape | Meaning |
|--------|---------|
| `\"` | Double quote |
| `\\` | Backslash |
| `\n` | Newline |
| `\t` | Tab |
| `\r` | Carriage return |

---

## Examples

### Complete Component

```pact
@note("Order management service")
component OrderService {
    type Order {
        +id: string
        +customer: Customer
        +items: OrderItem[]
        +status: OrderStatus
        -createdAt: string
    }

    type OrderItem {
        +product: Product
        +quantity: int
        +price: float
    }

    enum OrderStatus {
        Pending
        Processing
        Shipped
        Delivered
        Cancelled
    }

    depends on Database: database
    depends on PaymentGateway: external as payments
    depends on NotificationService

    provides OrderAPI {
        createOrder(customer: Customer, items: OrderItem[]) -> Order throws ValidationError
        getOrder(id: string) -> Order?
        async processOrder(id: string) -> Order throws PaymentError
    }

    requires InventoryService {
        checkStock(productId: string) -> int
        reserve(productId: string, quantity: int) -> bool
    }

    flow CreateOrder {
        validation = validateItems(items)
        if !validation.isValid {
            throw ValidationError
        }

        for item in items {
            stock = inventoryService.checkStock(item.product.id)
            if stock < item.quantity {
                throw InsufficientStockError
            }
        }

        order = db.createOrder(customer, items)
        await notifications.sendConfirmation(order)
        return order
    }

    states OrderLifecycle {
        initial Pending
        final Delivered
        final Cancelled

        state Processing {
            entry [validatePayment]
            exit [notifyCustomer]
        }

        Pending -> Processing on submit when payment.isValid()
        Processing -> Shipped on ship do [updateInventory, createShipment]
        Shipped -> Delivered on deliver
        Pending -> Cancelled on cancel
        Processing -> Cancelled on cancel when canCancel()
    }
}
```

---

Last updated: 2026-02-02
