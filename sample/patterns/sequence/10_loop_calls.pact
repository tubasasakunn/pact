// Pattern 10: Loop of calls (for)
component BatchProcessor {
    depends on ItemProcessor

    flow ProcessBatch {
        for item in items {
            validated = ItemProcessor.validate(item)
            processed = ItemProcessor.process(validated)
            saved = ItemProcessor.save(processed)
        }
        return true
    }
}

component ItemProcessor { }
