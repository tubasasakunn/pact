// Pattern: Sequential Flow with 3 steps
// 3ステップの順次処理フローパターン

component DataPipeline {
    type PipelineResult {
        recordsProcessed: int
        outputPath: string
    }

    flow ProcessData {
        // Step 1: Extract
        rawData = self.extractFromSource(sourceConfig)
        self.logExtractComplete(rawData)

        // Step 2: Transform
        transformedData = self.transformData(rawData)
        self.logTransformComplete(transformedData)

        // Step 3: Load
        outputPath = self.loadToDestination(transformedData)
        self.logLoadComplete(outputPath)

        return outputPath
    }
}
