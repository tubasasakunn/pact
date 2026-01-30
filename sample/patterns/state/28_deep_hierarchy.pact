// Pattern 28: 3-level nested states (deep hierarchy)
// Simulates hierarchical states with naming convention
component DeepHierarchy {
    states DocumentEditor {
        initial Idle
        final Closed

        // Top level states
        state Idle { }
        state Editing { }
        state Closed { }

        // Level 2: Editing substates (simulated with naming)
        state EditingText { }
        state EditingFormat { }
        state EditingSave { }

        // Level 3: EditingText substates
        state EditingTextTyping {
            entry [enableCursor, showKeyboard]
        }
        state EditingTextSelecting {
            entry [showSelectionHandles]
            exit [hideSelectionHandles]
        }
        state EditingTextCopying { }

        // Level 3: EditingFormat substates
        state EditingFormatFont { }
        state EditingFormatParagraph { }
        state EditingFormatStyle { }

        // Level 3: EditingSave substates
        state EditingSaveLocal {
            entry [saveToLocalStorage]
        }
        state EditingSaveCloud {
            entry [uploadToCloud]
            exit [notifyUploadComplete]
        }
        state EditingSaveExport { }

        // Top level transitions
        Idle -> Editing on openDocument
        Editing -> Closed on closeDocument

        // Level 2 transitions
        Editing -> EditingText on startTyping
        Editing -> EditingFormat on openFormatMenu
        Editing -> EditingSave on saveRequested

        EditingText -> EditingFormat on formatSelected
        EditingFormat -> EditingText on formatApplied
        EditingSave -> EditingText on saveComplete

        // Level 3 transitions within EditingText
        EditingText -> EditingTextTyping on keyPress
        EditingTextTyping -> EditingTextSelecting on shiftClick
        EditingTextSelecting -> EditingTextCopying on copyCommand
        EditingTextCopying -> EditingTextTyping on copyComplete
        EditingTextSelecting -> EditingTextTyping on escapePressed

        // Level 3 transitions within EditingFormat
        EditingFormat -> EditingFormatFont on fontMenu
        EditingFormat -> EditingFormatParagraph on paragraphMenu
        EditingFormat -> EditingFormatStyle on styleMenu
        EditingFormatFont -> EditingFormat on fontApplied
        EditingFormatParagraph -> EditingFormat on paragraphApplied
        EditingFormatStyle -> EditingFormat on styleApplied

        // Level 3 transitions within EditingSave
        EditingSave -> EditingSaveLocal on saveLocal
        EditingSave -> EditingSaveCloud on saveCloud
        EditingSave -> EditingSaveExport on exportFile
        EditingSaveLocal -> EditingSave on localSaved
        EditingSaveCloud -> EditingSave on cloudSaved
        EditingSaveExport -> EditingSave on exported

        // Cross-level transitions
        EditingTextTyping -> EditingSaveLocal on autoSave
        EditingFormatStyle -> Idle on discardChanges
    }
}
