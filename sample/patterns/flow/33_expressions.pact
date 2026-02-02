// Pattern 33: Various expression types
component ExpressionDemo {
    depends on Calculator : Calc as calc
    depends on DataService : Data as data

    flow TernaryExpressions {
        // Simple ternary
        result = condition ? valueIfTrue : valueIfFalse

        // Ternary with method calls
        status = isValid ? self.getSuccess() : self.getError()

        // Nested ternary
        grade = score > 90 ? gradeA : score > 80 ? gradeB : gradeC

        return result
    }

    flow NullCoalescing {
        // Simple null coalescing
        value = maybeNull ?? defaultValue

        // Null coalescing with throw
        required = maybeNull ?? throw RequiredError

        // Chained null coalescing
        result = first ?? second ?? third ?? defaultValue

        return result
    }

    flow BinaryExpressions {
        // Arithmetic
        sum = a + b
        diff = a - b
        product = a * b
        quotient = a / b
        remainder = a % b

        // Comparison
        isEqual = a == b
        isNotEqual = a != b
        isLess = a < b
        isGreater = a > b
        isLessOrEqual = a <= b
        isGreaterOrEqual = a >= b

        // Logical
        andResult = a && b
        orResult = a || b

        return sum
    }

    flow UnaryExpressions {
        // Logical not
        inverted = !flag

        // Negation
        negative = -value

        // Combined
        result = !isNegative ? value : -value

        return result
    }

    flow ComplexExpressions {
        // Combined operators
        result = (a + b) * c / d

        // Mixed with method calls
        computed = self.calculate(a + b) * data.getFactor()

        // Boolean expressions
        isValid = (a > 0) && (b < 100) || forceValid

        return computed
    }
}
