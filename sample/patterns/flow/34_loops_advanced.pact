// Pattern 34: Advanced loop patterns
component LoopPatterns {
    depends on DataStore : IStore as store
    depends on Processor : IProcessor as processor

    flow ForLoopVariants {
        // Simple for loop
        for item in items {
            processor.process(item)
        }

        // For loop with accumulator
        total = 0
        for value in values {
            total = total + value
        }

        // Nested for loops
        for row in rows {
            for cell in row.cells {
                processor.processCell(cell)
            }
        }

        // For loop with conditional
        for user in users {
            if user.isActive {
                processor.processActiveUser(user)
            }
        }

        return total
    }

    flow WhileLoopVariants {
        // Simple while loop
        while hasMore {
            item = store.getNext()
            processor.process(item)
        }

        // While with counter
        count = 0
        while count < maxIterations {
            processor.iterate(count)
            count = count + 1
        }

        // While with break condition
        while isRunning {
            result = processor.tryProcess()
            if result.isDone {
                isRunning = false
            }
        }

        return count
    }

    flow CombinedLoops {
        // For inside while
        while batchesRemaining {
            batch = store.getNextBatch()
            for item in batch.items {
                processor.processItem(item)
            }
        }

        // While inside for
        for queue in queues {
            while queue.hasMessages {
                message = queue.dequeue()
                processor.handleMessage(message)
            }
        }

        return true
    }

    flow LoopWithErrorHandling {
        errors = 0
        for item in items {
            result = processor.tryProcess(item)
            if result.isError {
                errors = errors + 1
                if errors > maxErrors {
                    throw TooManyErrorsException
                }
            }
        }
        return errors
    }

    flow PaginatedFetch {
        page = 0
        self.initResults()
        while hasMorePages {
            results = store.fetchPage(page)
            for result in results {
                self.appendResult(result)
            }
            page = page + 1
            hasMorePages = results.hasNext
        }
        allResults = self.getResults()
        return allResults
    }
}
