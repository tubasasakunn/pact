// Flow Pattern 04: Loop
component Service {
    flow ProcessItems {
        items = self.getItems()
        for item in items {
            self.processItem(item)
        }
        return true
    }
}
