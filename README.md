Easy command line utility to convert .TAP files (a data format for ZX-Spectrum emulator) into their WAV representation which canv be loaded by real device.

# Arguments
```
-a    amplify sound signal
-f int
      frequency of result wav, in Hz (default 22050)
-g int
      time gap between sound blocks, in seconds (default 1)
-i string
      source TAP file
-o string
      target WAV file
-s    add silence before the first file
```
# Example
```
zxtap2wav -i RENEGADE.tap
zxtap2wav -a -i RENEGADE.tap -o RENEGADE.wav -f 44100 -s
```