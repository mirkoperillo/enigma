This is only a fun project to improve my golang knowledge.

Enigma is a emulator of the [enigma](https://en.wikipedia.org/wiki/Enigma_machine) M3 machine, becamed famous during the II world war and used by nazist army to crypt
their communications

## Build

```
go build
```

## Quickstart

```
enigma [options] "the message"

Usage of ./enigma:
  -positions value
    	positions of rotors. Default A,A,A
  -reflector string
    	reflector rotor to use. Available values: B,C. Default B (default "B")
  -rotors value
    	comma separated list of rotors. Available values: I,II,III,IV,V. Default I,II,III
  -steckerboard value
    	comma separated list of pairs. Optional flag
  -verbose
    	more verbose output

```

An example:

```
./enigma -rotors=IV,I,III -positions=B,C,Q -steckerboard=A,Z,B,T -verbose "ciao github"
```

**NOTE: Following the historical usage practises, only letters are supported and the spaces are replaced with X letter**

## Resources

Some others emulators in the web:
* https://cryptii.com/pipes/enigma-machine
* https://www.101computing.net/enigma-machine-emulator/
