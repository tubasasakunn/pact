// Pattern 13: For loop
component ForLoopService {
    flow ForLoop {
        items = self.getItems()
        for item in items {
            processed = self.processItem(item)
            self.saveItem(processed)
        }
    }
}
