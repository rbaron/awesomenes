# 🎮 awesomenes

A basic NES emulator.

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

Games that use the [mapper](http://wiki.nesdev.com/w/index.php/Mapper) 0 mostly work, although without audio so far. Supporting more games (i.e. different mappers) are a priority.

# Roadmap

✅ CPU emulation

✅ Video support (picture processing unit - PPU)

✅ Keyboard input

✅ Mapper 0

➖  Joystick input

➖ More mappers

➖ Audio support (audio processing unit - PPU)


# Resources

All the information used to build this emulator was found on the awesome [nesdev wiki](https://wiki.nesdev.com).
