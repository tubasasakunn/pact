// Pattern: Sequential Flow with 4 steps
// 4ステップの順次処理フローパターン

component DeploymentPipeline {
    type DeployResult {
        version: string
        environment: string
        status: string
        url: string
    }

    flow Deploy {
        // Step 1: Build
        artifact = self.buildApplication(sourceCode)
        self.notifyBuildComplete(artifact)

        // Step 2: Test
        testResults = self.runTests(artifact)
        self.notifyTestsComplete(testResults)

        // Step 3: Stage
        stagingUrl = self.deployToStaging(artifact)
        self.runSmokeTests(stagingUrl)
        self.notifyStagingComplete(stagingUrl)

        // Step 4: Production
        productionUrl = self.deployToProduction(artifact)
        self.notifyProductionComplete(productionUrl)

        return productionUrl
    }
}
