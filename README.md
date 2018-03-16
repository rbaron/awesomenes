# 🎮 awesomenes

A NES emulator written in Go.

<p align="center">
  <img src="https://i.imgur.com/z8xYcxV.png" alt="dk"  width="400px"/>
  <img src="https://i.imgur.com/ahSN16z.png" alt="smb" width="400px"/>
</p>

# Getting and running

The easiest way to run `awesomenes` is to use the `go get` command:

```
$ go get github.com/rbaron/awesomenes
$ awesomenes MY_ROM.nes
```

# Status

Games that use the [mapper](http://wiki.nesdev.com/w/index.php/Mapper) 0 mostly work, although without audio so far. Supporting more games (i.e. different mappers) is a priority.

# Controller inputs

## Keyboard (controller 1)

```
Arrow keys  -> NES arrows
A           -> NES A
S           -> NES B
Enter       -> NES start
Right shift -> NES select
```

## Nintendo Switch Joycon (controller 1)

```
Directional -> NES arrows
Down arrow  -> NES A
Right arrow -> NES B
SL          -> NES select
SR          -> NES start
```

# Roadmap

✅ CPU emulation

✅ Video support (picture processing unit - PPU)

✅ Keyboard input

✅ Mapper 0

✅ Joystick input (tested with Nintendo Switch Joycon)

➖ More mappers

➖ Save state

➖ Audio support (audio processing unit - APU)


# Resources

All the information used to build this emulator was found on the awesome [nesdev wiki](https://wiki.nesdev.com).
