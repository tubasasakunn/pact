// Pattern 8: Multiple sequential calls
component DataPipeline {
    depends on Extractor
    depends on Transformer
    depends on Loader

    flow ETLProcess {
        raw = Extractor.extract(source)
        cleaned = Extractor.clean(raw)
        transformed = Transformer.transform(cleaned)
        validated = Transformer.validate(transformed)
        loaded = Loader.load(validated)
        verified = Loader.verify(loaded)
        return verified
    }
}

component Extractor { }
component Transformer { }
component Loader { }
