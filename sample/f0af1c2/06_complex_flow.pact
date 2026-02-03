// 06: Complex Flow - for/while loops, nested if, await, expressions
component DataPipeline {
    type Record {
        id: string
        data: string
        processed: bool
    }

    type BatchResult {
        total: int
        success: int
        failed: int
    }

    depends on DataSource
    depends on Transformer
    depends on DataSink
    depends on Monitor

    flow ProcessBatch {
        records = DataSource.fetchBatch(batchSize)
        successCount = 0
        failCount = 0

        for record in records {
            if record.data != null {
                transformed = Transformer.transform(record)
                if transformed.valid {
                    await DataSink.write(transformed)
                    successCount = successCount + 1
                } else {
                    Monitor.reportInvalid(record)
                    failCount = failCount + 1
                }
            } else {
                failCount = failCount + 1
            }
        }

        result = self.createResult(successCount, failCount)
        Monitor.reportBatch(result)
        return result
    }

    flow StreamProcess {
        running = true
        processed = 0

        while running {
            record = DataSource.nextRecord()
            if record != null {
                result = Transformer.transform(record)
                await DataSink.write(result)
                processed = processed + 1
            } else {
                running = false
            }
        }

        return processed
    }
}

component DataSource {
    provides SourceAPI {
        FetchBatch(size: int) -> string[]
        NextRecord() -> string
    }
}

component Transformer {
    provides TransformAPI {
        Transform(record: string) -> string
    }
}

component DataSink {
    provides SinkAPI {
        async Write(data: string) -> bool
    }
}

component Monitor {
    provides MonitorAPI {
        ReportInvalid(record: string)
        ReportBatch(result: string)
    }
}
