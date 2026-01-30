// Pattern 14: While loop
component WhileLoopService {
    flow WhileLoop {
        self.initialize()
        while hasMore {
            item = self.getNext()
            self.process(item)
        }
        self.finalize()
    }
}
