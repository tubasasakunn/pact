// Pattern 8: Simple if condition
component SimpleIfService {
    flow SimpleIf {
        value = self.getValue()
        if valid {
            self.processValid()
        }
    }
}
