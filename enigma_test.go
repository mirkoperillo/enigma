/*
      enigma - a M3 enigma machine emulator
      Written in 2021 by Mirko Perillo
      To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide.
	  This software is distributed without any warranty.
      You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.
*/

package main

import (
	"testing"
)

func TestEncodeA(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'A'}}
	encrypted := applyRotors(&cfg, 'A')
	if encrypted != 'Z' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeAWithPositionsAAB(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'B'}}
	encrypted := applyRotors(&cfg, 'A')
	if encrypted != 'F' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeAWithPositionsAAG(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'G'}}
	encrypted := applyRotors(&cfg, 'A')
	if encrypted != 'G' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeInverseRotorsAWithPositionsAAG(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'G'}, debug: false}
	encrypted := applyInverseRotors(&cfg, 'A')
	if encrypted != 'H' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeInverseRotorsLWithPositionsAAG(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'G'}, debug: false}
	encrypted := applyInverseRotors(&cfg, 'L')
	if encrypted != 'W' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeInverseRotorsAWithPositionsAAB(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'B'}, debug: false}
	encrypted := applyInverseRotors(&cfg, 'A')
	if encrypted != 'P' {
		t.Error("error: ", string(encrypted))
	}
}

func TestEncodeB(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'A'}}
	encrypted := applyRotors(&cfg, 'B')
	if encrypted != 'N' {
		t.Error("error: ", encrypted)
	}
}

func TestReverseRotor(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'A'}, debug: false}
	encrypted := applyInverseRotors(&cfg, 'A')
	if encrypted != 'D' {
		t.Error("error: ", encrypted)
	}
}

func TestReflector(t *testing.T) {
	cfg := config{reflector: reflectorB, debug: false}
	encrypted := applyReflector(&cfg, 'A')
	if encrypted != 'Y' {
		t.Error("error: ", encrypted)
	}
}

func TestSteckerBoard(t *testing.T) {
	cfg := config{steckerboard: [10]plug{plug{'A', 'Q'}, plug{'R', 'W'}}}
	encrypted := steckerboard(&cfg, 'Q')
	if encrypted != 'A' {
		t.Error("error: ", encrypted)
	}

}

func TestRotorsMoveInAAA(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'A'}}
	rotorsMovement(&cfg)
	newPositions := cfg.positions
	if differs(newPositions, [3]rune{'A', 'A', 'B'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}
}

func TestRotorsMoveInAAG(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'A', 'G'}}
	rotorsMovement(&cfg)
	newPositions := cfg.positions
	if differs(newPositions, [3]rune{'A', 'A', 'H'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}
}

func TestRotorsMoveInABC(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'B', 'C'}}
	rotorsMovement(&cfg)
	newPositions := cfg.positions
	if differs(newPositions, [3]rune{'A', 'B', 'D'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}
}

func TestRotorsMoveFirstNotch(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'B', 'V'}, debug: false}
	rotorsMovement(&cfg)
	newPositions := cfg.positions
	if differs(newPositions, [3]rune{'A', 'C', 'W'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}
}

func TestRotorsMoveDoubleStep(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, positions: [3]rune{'A', 'E', 'V'}}
	rotorsMovement(&cfg)
	newPositions := cfg.positions
	if differs(newPositions, [3]rune{'A', 'F', 'W'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}

	rotorsMovement(&cfg)
	newPositions = cfg.positions
	// double step
	if differs(newPositions, [3]rune{'B', 'G', 'X'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}

	rotorsMovement(&cfg)
	newPositions = cfg.positions
	// normale step
	if differs(newPositions, [3]rune{'B', 'G', 'Y'}) {
		t.Error("error positions: ", string(newPositions[0]), string(newPositions[1]), string(newPositions[2]))
	}
}

func TestScenario1(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'A', 'A'}, debug: false}
	msg := encrypt(&cfg, "enigma")
	if msg != "FQGAHW" {
		t.Error("error: ", msg)
	}
}

func TestScenario2(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'A', 'A'}, debug: false}
	msg := encrypt(&cfg, "enigmagoemulator")
	if msg != "FQGAHWZQWVGRBANF" {
		t.Error("error: ", msg)
	}
}

func TestScenarioWithSteckerboard(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'C', 'H'}, steckerboard: [10]plug{plug{'A', 'Q'}, plug{'R', 'W'}}, debug: false}
	msg := encrypt(&cfg, "gopherworld")
	if msg != "SFDXSDGJMJS" {
		t.Error("error: ", msg)
	}
}

func TestScenarioWithNoAlphabeticLetter(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'C', 'H'}, steckerboard: [10]plug{plug{'A', 'Q'}, plug{'R', 'W'}}, debug: false}
	msg := encrypt(&cfg, "gopher,world!!")
	if msg != "SFDXSDGJMJS" {
		t.Error("error: ", msg)
	}
}

func TestScenarioWithSpaces(t *testing.T) {
	cfg := config{rotors: [3]rotor{rotorI, rotorII, rotorIII}, reflector: reflectorB, positions: [3]rune{'A', 'C', 'H'}, steckerboard: [10]plug{plug{'A', 'Q'}, plug{'R', 'W'}}, debug: false}
	msg := encrypt(&cfg, "super gopher world")
	if msg != "GNDPGNWJBYUPGTQJGC" {
		t.Error("error: ", msg)
	}
}
func differs(a [3]rune, b [3]rune) bool {
	return a[0] != b[0] || a[1] != b[1] || a[2] != b[2]
}
