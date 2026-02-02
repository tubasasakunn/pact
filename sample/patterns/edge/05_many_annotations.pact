// Edge case: Component with 10+ annotations

@deprecated
@experimental
@internal
@public
@async
@synchronized
@cached
@logged
@traced
@monitored
@versioned("2.0")
@author("test")
component AnnotatedComponent {
    @nullable
    @readonly
    @indexed
    @unique
    @encrypted
    provides AnnotatedAPI {
        @deprecated
        @async
        @cached
        @logged
        @traced
        Process() -> Void
    }
}
