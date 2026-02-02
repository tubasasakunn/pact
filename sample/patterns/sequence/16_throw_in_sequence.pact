// Pattern 16: Throw in sequence
component ValidationService {
    depends on Validator

    flow ValidateData {
        result = Validator.check(data)
        if result == "invalid" {
            throw ValidationError
        }
        if result == "missing" {
            throw MissingDataError
        }
        if result == "duplicate" {
            throw DuplicateError
        }
        return result
    }
}

component Validator { }
