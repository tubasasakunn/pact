// Test Utilities - Helper functions for testing
// Source: internal/testutil/helper.go
@version("1.0")
component TestHelper {
	// Test context type
	type TestContext {
		testName: string
		failed: bool
		cleanupFuncs: func[]
	}

	depends on Parser

	// Test utility functions
	provides TestAPI {
		// Create a temporary directory for testing
		TempDir(t: TestContext) -> string

		// Write a file in the given directory
		WriteFile(t: TestContext, dir: string, name: string, content: string) -> string

		// Parse a string and fail test if parsing fails
		MustParse(t: TestContext, content: string) -> SpecFile

		// Assert string contains substring
		AssertContains(t: TestContext, str: string, substr: string)

		// Assert string does not contain substring
		AssertNotContains(t: TestContext, str: string, substr: string)

		// Assert no error occurred
		AssertNoError(t: TestContext, err: error)

		// Assert error occurred
		AssertError(t: TestContext, err: error)

		// Assert two values are equal
		AssertEqual(t: TestContext, got: any, want: any)

		// Truncate string to maximum length
		Truncate(str: string, maxLen: int) -> string

		// Get path to golden file in testdata/golden
		Golden(t: TestContext, name: string) -> string

		// Load golden file content
		LoadGolden(t: TestContext, name: string) -> string
	}

	// Create temp directory flow
	flow CreateTempDir {
		dir = self.mkTempDir("pact-test-*")
		if dir == null {
			self.failTest(t, "failed to create temp dir")
			throw TempDirError
		}
		self.registerCleanup(t, dir)
		return dir
	}

	// Write test file flow
	flow WriteTestFile {
		fullPath = self.joinPath(dir, name)
		dirPath = self.dirname(fullPath)
		mkdirResult = self.mkdirAll(dirPath)
		if mkdirResult == null {
			self.failTest(t, "failed to create dir")
			throw MkdirError
		}
		writeResult = self.writeFile(fullPath, content)
		if writeResult == null {
			self.failTest(t, "failed to write file")
			throw WriteError
		}
		return fullPath
	}

	// Assert contains flow
	flow CheckContains {
		isEmpty = self.isEmpty(input)
		if isEmpty {
			self.reportError(t, "string is empty")
			return
		}
		found = self.containsSubstring(input, expected)
		if found == false {
			self.reportError(t, "string does not contain expected substring")
		}
	}

	// Load golden file flow
	flow LoadGoldenFile {
		goldenPath = self.Golden(t, name)
		content = self.readFile(goldenPath)
		if content == null {
			self.failTest(t, "failed to load golden file " + name)
			throw GoldenFileError
		}
		return content
	}
}
