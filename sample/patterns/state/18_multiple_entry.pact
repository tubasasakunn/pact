// Pattern 18: State with multiple entry actions
component MultipleEntry {
    states Startup {
        initial Booting

        state Booting {
            entry [loadConfig, initHardware, startServices, logStartup]
        }
        state Ready { }

        Booting -> Ready on bootComplete
    }
}
