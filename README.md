[![License Apache 2.0](https://img.shields.io/badge/license-Apache%20License%202.0-green.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![PayPal donation](https://img.shields.io/badge/donation-PayPal-red.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=AHWJHJFBAWGL2)
[![Yandex.Money donation](https://img.shields.io/badge/donation-Я.деньги-yellow.svg)](http://yasobe.ru/na/iamoss)

Easy command line utility to convert [.TAP files (a data format for ZX-Spectrum emulator)](http://fileformats.archiveteam.org/wiki/TAP_(ZX_Spectrum)) into [sound WAV file](https://en.wikipedia.org/wiki/WAV).
(Also there is [similar utility to convert binary files into WAV files for personal computer BK-0010](https://github.com/raydac/bkbin2wav))

# How to build?
Just clone the project and use [maven](https://maven.apache.org/) `mvn package -Ppublish` command or [load needed prebuit version from the latest release](https://github.com/raydac/zxtap-to-wav/releases/tag/1.0.1).

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
# How to?

## Make longer silence interval between files in WAV
Just add `-g 2` or `-g 3` to make delay in 2 or 3 seconds.

## Add silence in start of generated WAV file
Use `-s` and silence will be generated in start of WAV file.

## I want 44100 Hz quantized WAV
Use parameter `-f 44100`

## Sound is too silent
Use flag `-a` and generated sound in WAV will be amplified to maximum.
