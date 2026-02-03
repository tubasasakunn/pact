// 10: Parallel States (Regions)
component MediaPlayer {
    type Track {
        title: string
        artist: string
        duration: int
    }

    states PlayerControl {
        initial Stopped
        final PoweredOff

        state Stopped {
            entry [resetPosition]
        }

        parallel Playing {
            region AudioRegion {
                initial Decoding

                state Decoding {
                    entry [startDecoder]
                }

                state Buffering {
                    entry [fillBuffer]
                }

                state Outputting {
                    entry [startAudioOutput]
                    exit [stopAudioOutput]
                }

                Decoding -> Buffering on decoded
                Buffering -> Outputting on bufferReady
                Outputting -> Decoding on nextFrame
            }

            region DisplayRegion {
                initial ShowingInfo

                state ShowingInfo {
                    entry [displayTrackInfo]
                }

                state ShowingLyrics {
                    entry [loadLyrics]
                }

                state ShowingVisualization {
                    entry [startVisualization]
                    exit [stopVisualization]
                }

                ShowingInfo -> ShowingLyrics on toggleLyrics
                ShowingLyrics -> ShowingVisualization on toggleViz
                ShowingVisualization -> ShowingInfo on toggleInfo
            }

            region VolumeRegion {
                initial Normal

                state Normal {
                    entry [setNormalVolume]
                }

                state Muted {
                    entry [muteAudio]
                    exit [unmuteAudio]
                }

                Normal -> Muted on mute
                Muted -> Normal on unmute
            }
        }

        Stopped -> Playing on play
        Playing -> Stopped on stop
        Stopped -> PoweredOff on powerOff
    }
}
