// Pact CLI - Command Line Interface Documentation
// This file documents the pact CLI tool architecture and behavior
@version("1.0.0")
@author("pact")
@description("Pact DSL diagram generator CLI")
component PactCLI {
    // =========================================================================
    // Types - Command Options
    // =========================================================================

    type GenerateOptions {
        output: string
        types: string[]
        files: string[]
    }

    type CheckOptions {
        showMissing: bool
        files: string[]
    }

    type CommandResult {
        success: bool
        message: string
        errorMessage: string
    }

    // =========================================================================
    // Types - Configuration
    // =========================================================================

    type PactConfig {
        src: string
        output: string
        types: string[]
    }

    type DiagramOutput {
        filename: string
        diagramType: string
        baseName: string
    }

    // =========================================================================
    // Types - Component Analysis
    // =========================================================================

    type ComponentInfo {
        name: string
        dependencies: string[]
    }

    type ValidationResult {
        file: string
        valid: bool
        errorMessage: string
    }

    // =========================================================================
    // Dependencies
    // =========================================================================

    depends on PactClient
    depends on FileSystem

    // =========================================================================
    // Provided Interface
    // =========================================================================

    provides CLI {
        Init() -> CommandResult
        Generate(opts: GenerateOptions) -> CommandResult
        Validate(files: string[]) -> CommandResult
        Check(opts: CheckOptions) -> CommandResult
        Watch(args: string[]) -> CommandResult
        Version() -> string
        Help() -> string
    }

    // =========================================================================
    // Command Execution States
    // =========================================================================

    states CommandExecution {
        initial Idle
        final Success
        final Error

        state Idle {
            entry [await_command]
        }

        state ParsingArgs {
            entry [parse_arguments]
        }

        state Executing {
            entry [start_timer]
            exit [stop_timer]
        }

        state Success {
            entry [print_success]
        }

        state Error {
            entry [print_error]
            exit [cleanup]
        }

        Idle -> ParsingArgs on command_received
        ParsingArgs -> Executing on args_valid
        ParsingArgs -> Error on args_invalid
        Executing -> Success on command_complete
        Executing -> Error on command_failed
    }

    // =========================================================================
    // Flow: Main Entry Point
    // =========================================================================

    flow Main {
        cmd = self.getCommand()
        cmdArgs = self.getArgs()
        result = self.routeCommand(cmd, cmdArgs)
        if errorOccurred {
            self.printError(result)
            return result
        }
        return result
    }

    // =========================================================================
    // Flow: Init Command
    // Creates default .pactconfig file
    // Usage: pact init
    // =========================================================================

    flow InitCommand {
        configPath = self.getConfigPath()
        exists = FileSystem.exists(configPath)
        if exists {
            self.printWarning("Config file already exists, overwriting...")
        }
        config = self.defaultConfig()
        FileSystem.writeFile(configPath, config)
        self.print("Created .pactconfig")
        return success
    }

    // =========================================================================
    // Flow: Generate Command
    // Generates diagrams from .pact files
    // Usage: pact generate [options] <files>
    // Options:
    //   -o, --output <dir>   Output directory (default: ".")
    //   -t, --type <types>   Diagram types: class,sequence,state,flow,all
    // =========================================================================

    flow GenerateCommand {
        opts = self.parseGenerateOptions(args)
        if noInputFiles {
            throw NoInputFilesError
        }
        outputDir = self.getOutputDir(opts)
        if createOutput {
            FileSystem.mkdirAll(outputDir)
        }
        files = self.expandFilePatterns(opts)
        if noFilesFound {
            throw NoPactFilesFoundError
        }
        for file in files {
            self.print("Processing " + file)
            spec = PactClient.parseFile(file)
            if parseError {
                throw ParseError
            }
            baseName = self.getBaseName(file)
            if generateClass {
                self.generateClassDiagram(spec, outputDir, baseName)
            }
            if generateSequence {
                self.generateSequenceDiagrams(spec, outputDir, baseName)
            }
            if generateState {
                self.generateStateDiagrams(spec, outputDir, baseName)
            }
            if generateFlow {
                self.generateFlowcharts(spec, outputDir, baseName)
            }
        }
        self.print("Done!")
        return success
    }

    // =========================================================================
    // Flow: Validate Command
    // Validates .pact files for syntax and semantic errors
    // Usage: pact validate <files>
    // =========================================================================

    flow ValidateCommand {
        if noFilesSpecified {
            throw NoFilesSpecifiedError
        }
        files = self.expandFilePatterns(args)
        hasErrors = false
        for file in files {
            result = PactClient.parseFile(file)
            if parseError {
                self.print("Error in " + file)
                hasErrors = true
            } else {
                self.print(file + ": OK")
            }
        }
        if hasErrors {
            throw ValidationFailedError
        }
        self.print("All files valid!")
        return success
    }

    // =========================================================================
    // Flow: Check Command
    // Checks for missing components and dependency issues
    // Usage: pact check [options] [files]
    // Options:
    //   -m, --missing   Show missing dependencies
    // =========================================================================

    flow CheckCommand {
        opts = self.parseCheckOptions(args)
        files = self.getFilesToCheck(opts)
        if noFilesFound {
            self.print("No .pact files found")
            if showMissing {
                throw NoFilesFoundError
            }
            return success
        }
        components = self.collectComponents(files)
        dependencies = self.collectDependencies(files)
        if showMissing {
            missing = self.findMissingDependencies(components, dependencies)
            if hasMissing {
                self.print("Missing components:")
                for name in missing {
                    self.print("  - " + name)
                }
                throw MissingComponentsError
            }
            self.print("No missing components")
        }
        self.print("Checked files and found components")
        return success
    }

    // =========================================================================
    // Flow: Watch Command
    // Watches for file changes and regenerates diagrams
    // Usage: pact watch [files]
    // Note: Not yet implemented
    // =========================================================================

    flow WatchCommand {
        self.print("Watch mode is not implemented yet")
        self.print("Use a file watcher like entr or fswatch with pact generate")
        return success
    }

    // =========================================================================
    // Helper Flow: Generate Class Diagram
    // =========================================================================

    flow GenerateClassDiagram {
        diagram = PactClient.toClassDiagram(spec)
        if diagramError {
            return error
        }
        outPath = self.buildOutputPath(output, baseName, "class")
        file = FileSystem.create(outPath)
        PactClient.renderClassDiagram(diagram, file)
        FileSystem.close(file)
        self.print("Generated " + baseName + "_class.svg")
        return success
    }

    // =========================================================================
    // Helper Flow: Generate Sequence Diagrams
    // =========================================================================

    flow GenerateSequenceDiagrams {
        flows = self.getFlowNames(spec)
        for flowName in flows {
            diagram = PactClient.toSequenceDiagram(spec, flowName)
            if diagramError {
                continue
            }
            outPath = self.buildOutputPath(output, baseName, "sequence", flowName)
            file = FileSystem.create(outPath)
            PactClient.renderSequenceDiagram(diagram, file)
            FileSystem.close(file)
            self.print("Generated sequence diagram for " + flowName)
        }
        return success
    }

    // =========================================================================
    // Helper Flow: Generate State Diagrams
    // =========================================================================

    flow GenerateStateDiagrams {
        stateNames = self.getStateNames(spec)
        for stateName in stateNames {
            diagram = PactClient.toStateDiagram(spec, stateName)
            if diagramError {
                continue
            }
            outPath = self.buildOutputPath(output, baseName, "state", stateName)
            file = FileSystem.create(outPath)
            PactClient.renderStateDiagram(diagram, file)
            FileSystem.close(file)
            self.print("Generated state diagram for " + stateName)
        }
        return success
    }

    // =========================================================================
    // Helper Flow: Generate Flowcharts
    // =========================================================================

    flow GenerateFlowcharts {
        flows = self.getFlowNames(spec)
        for flowName in flows {
            diagram = PactClient.toFlowchart(spec, flowName)
            if diagramError {
                continue
            }
            outPath = self.buildOutputPath(output, baseName, "flow", flowName)
            file = FileSystem.create(outPath)
            PactClient.renderFlowchart(diagram, file)
            FileSystem.close(file)
            self.print("Generated flowchart for " + flowName)
        }
        return success
    }
}

// =============================================================================
// External Dependency: Pact Client
// The core pact library for parsing and rendering
// =============================================================================

component PactClient {
    type SpecFile {
        hasComponent: bool
        componentCount: int
    }

    type ComponentData {
        name: string
        typeCount: int
        relationCount: int
        flowCount: int
        stateCount: int
    }

    type TypeDef {
        name: string
        fieldCount: int
    }

    type FieldDef {
        name: string
        fieldType: string
        optional: bool
    }

    type RelationDef {
        kind: string
        target: string
    }

    type FlowDef {
        name: string
        stepCount: int
    }

    type StatesDef {
        name: string
        stateCount: int
        transitionCount: int
    }

    type StateDef {
        name: string
        isInitial: bool
        isFinal: bool
        entryActions: string[]
        exitActions: string[]
    }

    type TransitionDef {
        fromState: string
        toState: string
        trigger: string
    }

    type InterfaceDef {
        name: string
        methodCount: int
    }

    type MethodDef {
        name: string
        paramCount: int
        hasReturn: bool
    }

    type ParamDef {
        name: string
        paramType: string
    }

    depends on FileSystem

    provides PactAPI {
        ParseFile(path: string) -> SpecFile
        ToClassDiagram(spec: SpecFile) -> string
        ToSequenceDiagram(spec: SpecFile, flowName: string) -> string
        ToStateDiagram(spec: SpecFile, stateName: string) -> string
        ToFlowchart(spec: SpecFile, flowName: string) -> string
        RenderClassDiagram(diagram: string, writer: string)
        RenderSequenceDiagram(diagram: string, writer: string)
        RenderStateDiagram(diagram: string, writer: string)
        RenderFlowchart(diagram: string, writer: string)
    }
}

// =============================================================================
// External Dependency: File System
// Standard file system operations
// =============================================================================

component FileSystem {
    type FileInfo {
        name: string
        isDir: bool
        size: int
    }

    type FileHandle {
        path: string
        handle: int
    }

    provides FileSystemAPI {
        Exists(path: string) -> bool
        Stat(path: string) -> FileInfo
        Create(path: string) -> FileHandle
        Close(file: FileHandle) -> bool
        WriteFile(path: string, content: string) -> bool
        ReadFile(path: string) -> string
        MkdirAll(path: string) -> bool
        Glob(pattern: string) -> string[]
    }
}
