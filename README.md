# wmc - WifiMCU CLI

This is the initial *alpha* release of a command-line "loader" for the [WifiMCU platform](http://www.wifimcu.com/).  This application is an alternative to some of the features offered in the Windows GUI app [WifiMCU Studio](https://github.com/SmartArduino/WiFiMCU-STUDIO)

## Platforms
* Linux (32-bit, 64-bit, arm)
* Windows (32-bit, 64-bit)
* Mac (64-bit)

## Usage
```
  wmc [command]

Available Commands:
  ver         Get the current version
  ls          List all files on device
  put         Send a file to the device
  rm          Remove a file
  config      Display current config
  cmd         Run an arbitrary command, needs to be in double-quotes
  read        read a file
  help        Help about any command

Flags:
  -h, --help=false: help for wmc


Use "wmc [command] --help" for more information about a command.
```

## Get Involved
Please submit any bugs or feature request to the Github issue tracker
