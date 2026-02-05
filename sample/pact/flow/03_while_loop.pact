// Pattern: While Loop Flow
// ループを含むフローパターン

component BatchProcessor {
    type BatchResult {
        processed: int
        failed: int
        duration: float
    }

    flow ProcessBatch {
        items = self.getQueuedItems()
        processed = 0
        failed = 0

        while items.hasNext {
            item = items.next
            result = self.processItem(item)

            if result.success {
                processed = self.incrementProcessed(processed)
                self.markComplete(item)
            } else {
                failed = self.incrementFailed(failed)
                self.markFailed(item)
            }
        }

        duration = self.calculateDuration()
        return duration
    }
}
