// Pattern 6: Call with return value
component OrderService {
    depends on PriceCalculator

    flow GetTotal {
        price = PriceCalculator.calculate(items)
        tax = PriceCalculator.getTax(price)
        total = PriceCalculator.addTax(price, tax)
        return total
    }
}

component PriceCalculator {
    provides PricingAPI {
        Calculate(items: string[]) -> float
        GetTax(amount: float) -> float
        AddTax(amount: float, tax: float) -> float
    }
}
