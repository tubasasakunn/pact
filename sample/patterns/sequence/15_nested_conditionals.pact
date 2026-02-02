// Pattern 15: Nested conditionals
component ShippingService {
    depends on LocalCarrier
    depends on InternationalCarrier
    depends on ExpressCarrier

    flow SelectCarrier {
        if destination == "domestic" {
            if priority == "express" {
                rate = ExpressCarrier.getRate(package)
            } else {
                rate = LocalCarrier.getRate(package)
            }
        } else {
            if size == "oversized" {
                rate = InternationalCarrier.getOversizedRate(package)
            } else {
                rate = InternationalCarrier.getStandardRate(package)
            }
        }
        return rate
    }
}

component LocalCarrier { }
component InternationalCarrier { }
component ExpressCarrier { }
