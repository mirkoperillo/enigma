/*
enigma - a M3 enigma machine emulator
Written in 2021 by Mirko Perillo
To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide.
This software is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.
*/
package enigma

import (
	"fmt"
	"strings"
)

const LETTERS_IN_ALPHABET = 26

var RotorI = Rotor{[26]rune{'E', 'K', 'M', 'F', 'L', 'G', 'D', 'Q', 'V', 'Z', 'N', 'T', 'O', 'W', 'Y', 'H', 'X', 'U', 'S', 'P', 'A', 'I', 'B', 'R', 'C', 'J'}, 'R'}
var RotorII = Rotor{[26]rune{'A', 'J', 'D', 'K', 'S', 'I', 'R', 'U', 'X', 'B', 'L', 'H', 'W', 'T', 'M', 'C', 'Q', 'G', 'Z', 'N', 'P', 'Y', 'F', 'V', 'O', 'E'}, 'F'}
var RotorIII = Rotor{[26]rune{'B', 'D', 'F', 'H', 'J', 'L', 'C', 'P', 'R', 'T', 'X', 'V', 'Z', 'N', 'Y', 'E', 'I', 'W', 'G', 'A', 'K', 'M', 'U', 'S', 'Q', 'O'}, 'W'}
var RotorIV = Rotor{[26]rune{'E', 'S', 'O', 'V', 'P', 'Z', 'J', 'A', 'Y', 'Q', 'U', 'I', 'R', 'H', 'X', 'L', 'N', 'F', 'T', 'G', 'K', 'D', 'C', 'M', 'W'}, 'K'}
var RotorV = Rotor{[26]rune{'V', 'Z', 'B', 'R', 'G', 'I', 'T', 'Y', 'U', 'P', 'S', 'D', 'N', 'H', 'L', 'X', 'A', 'W', 'M', 'J', 'Q', 'O', 'F', 'E', 'C', 'K'}, 'A'}
var ReflectorB = Rotor{[26]rune{'Y', 'R', 'U', 'H', 'Q', 'S', 'L', 'D', 'P', 'X', 'N', 'G', 'O', 'K', 'M', 'I', 'E', 'B', 'F', 'Z', 'C', 'W', 'V', 'J', 'A', 'T'}, '0'}
var ReflectorC = Rotor{[26]rune{'F', 'V', 'P', 'J', 'I', 'A', 'O', 'Y', 'E', 'D', 'R', 'Z', 'X', 'W', 'G', 'C', 'T', 'K', 'U', 'Q', 'S', 'B', 'N', 'M', 'H', 'L'}, '0'}

type Plug struct {
	A rune
	B rune
}

type Rotor struct {
	Letters [26]rune
	Notch   rune
}
type Config struct {
	Steckerboard     [10]Plug
	Rotors           [3]Rotor
	Reflector        Rotor
	Positions        [3]rune
	DoubleRotorsStep bool
	Rings            [3]int
	Debug            bool
}

func applyRotors(cfg *Config, letter rune) rune {
	position := letter - 'A'
	var encryptedLetter rune
	for i := len(cfg.Rotors) - 1; i >= 0; i-- {
		rotorPosition := cfg.Positions[i] - 'A'
		encryptedLetter = cfg.Rotors[i].Letters[((position + rotorPosition) % LETTERS_IN_ALPHABET)]
		encryptedLetter = ((LETTERS_IN_ALPHABET + encryptedLetter - 'A' - rotorPosition) % LETTERS_IN_ALPHABET) + 'A'
		position = encryptedLetter - 'A'
		if cfg.Debug {
			fmt.Println("Rotor step: ", string(encryptedLetter))
		}
	}
	return encryptedLetter
}

func applyInverseRotors(cfg *Config, letter rune) rune {
	var encryptedLetter rune
	var position int
	for i, r := range cfg.Rotors {
		rotorPosition := cfg.Positions[i] - 'A'
		letter = 'A' + ((letter - 'A' + rotorPosition) % LETTERS_IN_ALPHABET)
		for pos, l := range r.Letters {
			if l == letter {
				position = pos
				break
			}
		}
		encryptedLetter = 'A' + ((LETTERS_IN_ALPHABET + rune(position) - rotorPosition) % LETTERS_IN_ALPHABET)
		letter = encryptedLetter
		if cfg.Debug {
			fmt.Println("Inverse rotor step: ", string(encryptedLetter))
		}
	}
	return encryptedLetter
}

func applyReflector(cfg *Config, letter rune) rune {
	position := letter - 'A'
	return cfg.Reflector.Letters[position]
}

func steckerboard(cfg *Config, letter rune) rune {
	encrypted := letter
	for _, plug := range cfg.Steckerboard {
		if plug.A == letter {
			encrypted = plug.B
			break
		} else if plug.B == letter {
			encrypted = plug.A
			break
		}
	}
	return encrypted
}

func rotorsMovement(cfg *Config) {
	if cfg.DoubleRotorsStep {
		currentPos := cfg.Positions[1]
		cfg.Positions[1] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		currentPos = cfg.Positions[0]
		cfg.Positions[0] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		cfg.DoubleRotorsStep = false
		if cfg.Debug {
			fmt.Println("doubleRotorStep executed")
		}
	}
	currentPos := cfg.Positions[2]
	nextPos := 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
	cfg.Positions[2] = nextPos
	if cfg.Positions[2] == cfg.Rotors[2].Notch {
		if cfg.Debug {
			fmt.Println("Rotor position 3 notch")
		}
		currentPos = cfg.Positions[1]
		cfg.Positions[1] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		if cfg.Positions[1] == cfg.Rotors[1].Notch {
			cfg.DoubleRotorsStep = true
			if cfg.Debug {
				fmt.Println("Rotor position 2 notch, doubleRotorStep on next movement")
			}
		}
	}
}

func Encrypt(cfg *Config, msg string) string {
	msg = strings.ToUpper(msg)
	if cfg.Debug {
		fmt.Println("Input: ", msg)
	}
	var encryptedMsg string

	for _, l := range msg {

		if l == ' ' {
			l = 'X'
		}

		if l >= 65 && l <= 90 {
			// 0. rotors movement
			rotorsMovement(cfg)
			if cfg.Debug {
				fmt.Println("positions: ", string(cfg.Positions[0]), string(cfg.Positions[1]), string(cfg.Positions[2]))
			}
			if cfg.Debug {
				fmt.Println("Input step: ", string(l))
			}
			// 1. stecker
			encryptedLetter := steckerboard(cfg, l)
			if cfg.Debug {
				fmt.Println("steckerboard encryption: ", string(encryptedLetter))
			}

			// 2. rotors
			encryptedLetter = applyRotors(cfg, encryptedLetter)
			if cfg.Debug {
				fmt.Println("rotors encryption: ", string(encryptedLetter))
			}
			// 3. reflector
			encryptedLetter = applyReflector(cfg, encryptedLetter)
			if cfg.Debug {
				fmt.Println("reflector encryption: ", string(encryptedLetter))
			}
			// 4. inverse rotors
			encryptedLetter = applyInverseRotors(cfg, encryptedLetter)
			if cfg.Debug {
				fmt.Println("rotors encryption: ", string(encryptedLetter))
			}
			// 5. stecker
			encryptedLetter = steckerboard(cfg, encryptedLetter)
			if cfg.Debug {
				fmt.Println("steckerboard encryption: ", string(encryptedLetter))
			}

			encryptedMsg += string(encryptedLetter)
		}
	}
	return encryptedMsg
}
