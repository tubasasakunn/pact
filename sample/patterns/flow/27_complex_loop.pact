// Pattern 27: Nested for and while loops with complex iteration
component ComplexLoopService {
    depends on BatchProcessor
    depends on ItemService

    flow ComplexLoop {
        batches = self.getBatches()
        processedCount = self.initCounter()

        // Outer for loop over batches
        for batch in batches {
            items = BatchProcessor.getItems(batch)
            batchResult = self.initBatchResult(batch)

            // Inner for loop over items in batch
            for item in items {
                validated = ItemService.validate(item)
                if itemValid {
                    processed = ItemService.process(item)
                    self.addToResult(batchResult, processed)
                } else {
                    self.logInvalidItem(item)
                }
            }

            // While loop for retry logic
            retryCount = self.getRetryCount()
            while hasRetries {
                failedItems = self.getFailedItems(batchResult)
                for failedItem in failedItems {
                    retryResult = ItemService.retry(failedItem)
                    if retrySucceeded {
                        self.updateResult(batchResult, retryResult)
                    }
                }
                retryCount = self.decrementRetry(retryCount)
            }

            BatchProcessor.saveBatchResult(batchResult)
            processedCount = self.incrementCounter(processedCount)
        }

        summary = self.generateSummary(processedCount)
        return summary
    }
}
