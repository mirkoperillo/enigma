/*
   enigma - a M3 enigma machine emulator
   Written in 2021 by Mirko Perillo
   To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide.
   This software is distributed without any warranty.
   You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const LETTERS_IN_ALPHABET = 26

var rotorI = rotor{[26]rune{'E', 'K', 'M', 'F', 'L', 'G', 'D', 'Q', 'V', 'Z', 'N', 'T', 'O', 'W', 'Y', 'H', 'X', 'U', 'S', 'P', 'A', 'I', 'B', 'R', 'C', 'J'}, 'R'}
var rotorII = rotor{[26]rune{'A', 'J', 'D', 'K', 'S', 'I', 'R', 'U', 'X', 'B', 'L', 'H', 'W', 'T', 'M', 'C', 'Q', 'G', 'Z', 'N', 'P', 'Y', 'F', 'V', 'O', 'E'}, 'F'}
var rotorIII = rotor{[26]rune{'B', 'D', 'F', 'H', 'J', 'L', 'C', 'P', 'R', 'T', 'X', 'V', 'Z', 'N', 'Y', 'E', 'I', 'W', 'G', 'A', 'K', 'M', 'U', 'S', 'Q', 'O'}, 'W'}
var rotorIV = rotor{[26]rune{'E', 'S', 'O', 'V', 'P', 'Z', 'J', 'A', 'Y', 'Q', 'U', 'I', 'R', 'H', 'X', 'L', 'N', 'F', 'T', 'G', 'K', 'D', 'C', 'M', 'W'}, 'K'}
var rotorV = rotor{[26]rune{'V', 'Z', 'B', 'R', 'G', 'I', 'T', 'Y', 'U', 'P', 'S', 'D', 'N', 'H', 'L', 'X', 'A', 'W', 'M', 'J', 'Q', 'O', 'F', 'E', 'C', 'K'}, 'A'}
var reflectorB = rotor{[26]rune{'Y', 'R', 'U', 'H', 'Q', 'S', 'L', 'D', 'P', 'X', 'N', 'G', 'O', 'K', 'M', 'I', 'E', 'B', 'F', 'Z', 'C', 'W', 'V', 'J', 'A', 'T'}, '0'}
var reflectorC = rotor{[26]rune{'F', 'V', 'P', 'J', 'I', 'A', 'O', 'Y', 'E', 'D', 'R', 'Z', 'X', 'W', 'G', 'C', 'T', 'K', 'U', 'Q', 'S', 'B', 'N', 'M', 'H', 'L'}, '0'}

type plug struct {
	a rune
	b rune
}

type rotor struct {
	letters [26]rune
	notch   rune
}
type config struct {
	steckerboard     [10]plug
	rotors           [3]rotor
	reflector        rotor
	positions        [3]rune
	doubleRotorsStep bool
	rings            [3]int
	debug            bool
}

type arrayFlags []string

func (a *arrayFlags) Set(value string) error {
	values := strings.Split(value, ",")
	for _, v := range values {
		*a = append(*a, strings.ToUpper(v))
	}
	return nil
}

func (a *arrayFlags) String() string {
	return strings.Join(*a, ",")
}

type steckerFlags []plug

func (s *steckerFlags) Set(value string) error {
	values := strings.Split(value, ",")
	if len(values)%2 != 0 {
		return errors.New("Steckerboard options should be even")
	}
	for i := 0; i < len(values); i = i + 2 {
		a := string2rune(values[i])
		b := string2rune(values[i+1])
		*s = append(*s, plug{a, b})
	}

	return nil
}

func (s *steckerFlags) String() string {
	var result string
	for _, pair := range *s {
		result = result + string(pair.a) + "," + string(pair.b)
	}
	return result
}

func string2rune(s string) rune {
	r := []rune(strings.ToUpper(s))
	return r[0]
}

func validateRotorsFlag(rotors []string) (bool, error) {
	var validNames = [5]string{"I", "II", "III", "IV", "V"}
	if len(rotors) != 3 {
		panic("You have to select 3 rotors")
	}

	for _, r := range rotors {
		if !isInArray(r, validNames) {
			return false, errors.New(fmt.Sprintf("Rotor %s not exist", r))
		}
	}
	return true, nil
}

func validateReflectorFlag(reflector string) (bool, error) {
	var validNames = [5]string{"B", "C"}
	if !isInArray(reflector, validNames) {
		return false, errors.New(fmt.Sprintf("Reflector %s not exist", reflector))
	}
	return true, nil
}

func validatePositionsFlag(positions []string) (bool, error) {
	if len(positions) != 3 {
		panic("You have to select 3 positions")
	}
	for _, r := range positions {
		letter := []rune(r)
		if !isLetter(letter[0]) {
			return false, errors.New(fmt.Sprintf("Position %s is not valid", string(letter[0])))
		}
	}
	return true, nil
}

func isLetter(r rune) bool {
	return r >= 65 && r <= 90
}

func validateSteckerboardFlag(stecker []plug) (bool, error) {
	for _, p := range stecker {
		if !isLetter(p.a) || !isLetter(p.b) {
			return false, errors.New(fmt.Sprintf("Not valid value in (%s,%s)", string(p.a), string(p.b)))
		}
	}
	return true, nil
}

func isInArray(elem string, collection [5]string) bool {
	for _, e := range collection {
		if e == elem {
			return true
		}
	}
	return false
}

func flagToRotor(flag string) rotor {
	switch flag {
	case "I":
		return rotorI
	case "II":
		return rotorII
	case "III":
		return rotorIII
	case "IV":
		return rotorIV
	case "V":
		return rotorV
	case "B":
		return reflectorB
	case "C":
		return reflectorC
	default:
		return rotorI
	}
}
func toConfig(rotors arrayFlags, reflector string, positions arrayFlags, steckerboard steckerFlags, verbose bool) config {
	debugCfg := verbose
	var steckerboardCfg [10]plug
	for idx, p := range steckerboard {
		steckerboardCfg[idx] = p
	}

	reflectorCfg := flagToRotor(reflector)
	rotorsCfg := [3]rotor{flagToRotor(rotors[0]), flagToRotor(rotors[1]), flagToRotor(rotors[2])}
	positionsCfg := [3]rune{string2rune(positions[0]), string2rune(positions[1]), string2rune(positions[2])}
	return config{rotors: rotorsCfg, reflector: reflectorCfg, positions: positionsCfg, steckerboard: steckerboardCfg, debug: debugCfg}
}

func main() {
	// example of commang flags
	// --rotors=I,II,IV --reflector=B --positions=A,A,A --verbose --steckerboard=A,B,C,E
	lastArg := os.Args[len(os.Args)-1]
	var message string
	if strings.HasPrefix(lastArg, "--") {
		message = ""
	} else {
		message = lastArg
	}
	var rotorsFlag arrayFlags
	var reflectorFlag *string
	var positionsFlag arrayFlags
	var steckerboardFlag steckerFlags
	var verboseFlag *bool

	flag.Var(&rotorsFlag, "rotors", "comma separated list of rotors. Available values: I,II,III,IV,V. Default I,II,III")
	reflectorFlag = flag.String("reflector", "B", "reflector rotor to use. Available values: B,C. Default B")
	flag.Var(&positionsFlag, "positions", "positions of rotors. Default A,A,A")
	flag.Var(&steckerboardFlag, "steckerboard", "comma separated list of pairs. Optional flag")
	verboseFlag = flag.Bool("verbose", false, "more verbose output")

	flag.Parse()
	if rotorsFlag == nil {
		rotorsFlag = []string{"I", "II", "III"}
	} else {
		_, err := validateRotorsFlag(rotorsFlag)
		if err != nil {
			panic(err)
		}
	}
	_, err := validateReflectorFlag(*reflectorFlag)
	if err != nil {
		panic(err)
	}
	if positionsFlag == nil {
		positionsFlag = []string{"A", "A", "A"}
	} else {
		_, err := validatePositionsFlag(positionsFlag)
		if err != nil {
			panic(err)
		}
	}
	_, err = validateSteckerboardFlag(steckerboardFlag)
	if err != nil {
		panic(err)
	}

	config := toConfig(rotorsFlag, *reflectorFlag, positionsFlag, steckerboardFlag, *verboseFlag)
	resultMsg := encrypt(&config, message)
	fmt.Println(resultMsg)
}

func applyRotors(cfg *config, letter rune) rune {
	position := letter - 'A'
	var encryptedLetter rune
	for i := len(cfg.rotors) - 1; i >= 0; i-- {
		rotorPosition := cfg.positions[i] - 'A'
		encryptedLetter = cfg.rotors[i].letters[((position + rotorPosition) % LETTERS_IN_ALPHABET)]
		encryptedLetter = ((LETTERS_IN_ALPHABET + encryptedLetter - 'A' - rotorPosition) % LETTERS_IN_ALPHABET) + 'A'
		position = encryptedLetter - 'A'
		if cfg.debug {
			fmt.Println("Rotor step: ", string(encryptedLetter))
		}
	}
	return encryptedLetter
}

func applyInverseRotors(cfg *config, letter rune) rune {
	var encryptedLetter rune
	var position int
	for i, r := range cfg.rotors {
		rotorPosition := cfg.positions[i] - 'A'
		letter = 'A' + ((letter - 'A' + rotorPosition) % LETTERS_IN_ALPHABET)
		for pos, l := range r.letters {
			if l == letter {
				position = pos
				break
			}
		}
		encryptedLetter = 'A' + ((LETTERS_IN_ALPHABET + rune(position) - rotorPosition) % LETTERS_IN_ALPHABET)
		letter = encryptedLetter
		if cfg.debug {
			fmt.Println("Inverse rotor step: ", string(encryptedLetter))
		}
	}
	return encryptedLetter
}

func applyReflector(cfg *config, letter rune) rune {
	position := letter - 'A'
	return cfg.reflector.letters[position]
}

func steckerboard(cfg *config, letter rune) rune {
	encrypted := letter
	for _, plug := range cfg.steckerboard {
		if plug.a == letter {
			encrypted = plug.b
			break
		} else if plug.b == letter {
			encrypted = plug.a
			break
		}
	}
	return encrypted
}

func rotorsMovement(cfg *config) {
	if cfg.doubleRotorsStep {
		currentPos := cfg.positions[1]
		cfg.positions[1] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		currentPos = cfg.positions[0]
		cfg.positions[0] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		cfg.doubleRotorsStep = false
		if cfg.debug {
			fmt.Println("doubleRotorStep executed")
		}
	}
	currentPos := cfg.positions[2]
	nextPos := 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
	cfg.positions[2] = nextPos
	if cfg.positions[2] == cfg.rotors[2].notch {
		if cfg.debug {
			fmt.Println("Rotor position 3 notch")
		}
		currentPos = cfg.positions[1]
		cfg.positions[1] = 'A' + (currentPos+1-'A')%LETTERS_IN_ALPHABET
		if cfg.positions[1] == cfg.rotors[1].notch {
			cfg.doubleRotorsStep = true
			if cfg.debug {
				fmt.Println("Rotor position 2 notch, doubleRotorStep on next movement")
			}
		}
	}
}

func encrypt(cfg *config, msg string) string {
	msg = strings.ToUpper(msg)
	if cfg.debug {
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
			if cfg.debug {
				fmt.Println("positions: ", string(cfg.positions[0]), string(cfg.positions[1]), string(cfg.positions[2]))
			}
			if cfg.debug {
				fmt.Println("Input step: ", string(l))
			}
			// 1. stecker
			encryptedLetter := steckerboard(cfg, l)
			if cfg.debug {
				fmt.Println("steckerboard encryption: ", string(encryptedLetter))
			}

			// 2. rotors
			encryptedLetter = applyRotors(cfg, encryptedLetter)
			if cfg.debug {
				fmt.Println("rotors encryption: ", string(encryptedLetter))
			}
			// 3. reflector
			encryptedLetter = applyReflector(cfg, encryptedLetter)
			if cfg.debug {
				fmt.Println("reflector encryption: ", string(encryptedLetter))
			}
			// 4. inverse rotors
			encryptedLetter = applyInverseRotors(cfg, encryptedLetter)
			if cfg.debug {
				fmt.Println("rotors encryption: ", string(encryptedLetter))
			}
			// 5. stecker
			encryptedLetter = steckerboard(cfg, encryptedLetter)
			if cfg.debug {
				fmt.Println("steckerboard encryption: ", string(encryptedLetter))
			}

			encryptedMsg += string(encryptedLetter)
		}
	}
	return encryptedMsg
}
