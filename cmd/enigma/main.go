package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mirkoperillo/enigma/internal/enigma"
)

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

type steckerFlags []enigma.Plug

func (s *steckerFlags) Set(value string) error {
	values := strings.Split(value, ",")
	if len(values)%2 != 0 {
		return errors.New("Steckerboard options should be even")
	}
	for i := 0; i < len(values); i = i + 2 {
		a := string2rune(values[i])
		b := string2rune(values[i+1])
		*s = append(*s, enigma.Plug{A: a, B: b})
	}

	return nil
}

func string2rune(s string) rune {
	r := []rune(strings.ToUpper(s))
	return r[0]
}

func (s *steckerFlags) String() string {
	var result string
	for _, pair := range *s {
		result = result + string(pair.A) + "," + string(pair.B)
	}
	return result
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
	resultMsg := enigma.Encrypt(&config, message)
	fmt.Println(resultMsg)
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

func toConfig(rotors arrayFlags, reflector string, positions arrayFlags, steckerboard steckerFlags, verbose bool) enigma.Config {
	debugCfg := verbose
	var steckerboardCfg [10]enigma.Plug
	for idx, p := range steckerboard {
		steckerboardCfg[idx] = p
	}

	reflectorCfg := flagToRotor(reflector)
	rotorsCfg := [3]enigma.Rotor{flagToRotor(rotors[0]), flagToRotor(rotors[1]), flagToRotor(rotors[2])}
	positionsCfg := [3]rune{string2rune(positions[0]), string2rune(positions[1]), string2rune(positions[2])}
	return enigma.Config{Rotors: rotorsCfg, Reflector: reflectorCfg, Positions: positionsCfg, Steckerboard: steckerboardCfg, Debug: debugCfg}
}

func flagToRotor(flag string) enigma.Rotor {
	switch flag {
	case "I":
		return enigma.RotorI
	case "II":
		return enigma.RotorII
	case "III":
		return enigma.RotorIII
	case "IV":
		return enigma.RotorIV
	case "V":
		return enigma.RotorV
	case "B":
		return enigma.ReflectorB
	case "C":
		return enigma.ReflectorC
	default:
		return enigma.RotorI
	}
}

func validateSteckerboardFlag(stecker []enigma.Plug) (bool, error) {
	for _, p := range stecker {
		if !isLetter(p.A) || !isLetter(p.B) {
			return false, errors.New(fmt.Sprintf("Not valid value in (%s,%s)", string(p.A), string(p.B)))
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

func isLetter(r rune) bool {
	return r >= 65 && r <= 90
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

func validateReflectorFlag(reflector string) (bool, error) {
	var validNames = [5]string{"B", "C"}
	if !isInArray(reflector, validNames) {
		return false, errors.New(fmt.Sprintf("Reflector %s not exist", reflector))
	}
	return true, nil
}
