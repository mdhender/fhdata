// fhdata - Far Horizons Data
//
// Copyright (c) 2022 Michael D Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package fhdata

func codeToGas(code int) Gas {
	switch code {
	case 1:
		return Gas{Code: "H2", Name: "Hydrogen"}
	case 2:
		return Gas{Code: "CH4", Name: "Methane"}
	case 3:
		return Gas{Code: "He", Name: "Helium"}
	case 4:
		return Gas{Code: "NH3", Name: "Ammonia"}
	case 5:
		return Gas{Code: "N2", Name: "Nitrogen"}
	case 6:
		return Gas{Code: "CO2", Name: "Carbon Dioxide"}
	case 7:
		return Gas{Code: "O2", Name: "Oxygen"}
	case 8:
		return Gas{Code: "HCl", Name: "Hydrogen Chloride"}
	case 9:
		return Gas{Code: "Cl2", Name: "Chlorine"}
	case 10:
		return Gas{Code: "F2", Name: "Fluorine"}
	case 11:
		return Gas{Code: "H2O", Name: "Steam"}
	case 12:
		return Gas{Code: "SO2", Name: "Sulfur Dioxide"}
	case 13:
		return Gas{Code: "H2S", Name: "Hydrogen Sulfide"}
	default:
		return Gas{}
	}
}

func codeToItem(code int, qty int) Item {
	switch code {
	case 0:
		return Item{Code: "RM", Name: "Raw Material Unit", Quantity: qty}
	case 1:
		return Item{Code: "PD", Name: "Planetary Defense Unit", Quantity: qty}
	case 2:
		return Item{Code: "SU", Name: "Starbase Unit", Quantity: qty}
	case 3:
		return Item{Code: "DR", Name: "Damage Repair Unit", Quantity: qty}
	case 4:
		return Item{Code: "CU", Name: "Colonist Unit", Quantity: qty}
	case 5:
		return Item{Code: "IU", Name: "Colonial Mining Unit", Quantity: qty}
	case 6:
		return Item{Code: "AU", Name: "Colonial Manufacturing Unit", Quantity: qty}
	case 7:
		return Item{Code: "FS", Name: "Fail-Safe Jump Unit", Quantity: qty}
	case 8:
		return Item{Code: "JP", Name: "Jump Portal Unit", Quantity: qty}
	case 9:
		return Item{Code: "FM", Name: "Forced Misjump Unit", Quantity: qty}
	case 10:
		return Item{Code: "FJ", Name: "Forced Jump Unit", Quantity: qty}
	case 11:
		return Item{Code: "GT", Name: "Gravitic Telescope Unit", Quantity: qty}
	case 12:
		return Item{Code: "FD", Name: "Field Distortion Unit", Quantity: qty}
	case 13:
		return Item{Code: "TP", Name: "Terraforming Plant", Quantity: qty}
	case 14:
		return Item{Code: "GW", Name: "Germ Warfare Bomb", Quantity: qty}
	case 15:
		return Item{Code: "SG1", Name: "Mark-1 Shield Generator", Quantity: qty}
	case 16:
		return Item{Code: "SG2", Name: "Mark-2 Shield Generator", Quantity: qty}
	case 17:
		return Item{Code: "SG3", Name: "Mark-3 Shield Generator", Quantity: qty}
	case 18:
		return Item{Code: "SG4", Name: "Mark-4 Shield Generator", Quantity: qty}
	case 19:
		return Item{Code: "SG5", Name: "Mark-5 Shield Generator", Quantity: qty}
	case 20:
		return Item{Code: "SG6", Name: "Mark-6 Shield Generator", Quantity: qty}
	case 21:
		return Item{Code: "SG7", Name: "Mark-7 Shield Generator", Quantity: qty}
	case 22:
		return Item{Code: "SG8", Name: "Mark-8 Shield Generator", Quantity: qty}
	case 23:
		return Item{Code: "SG9", Name: "Mark-9 Shield Generator", Quantity: qty}
	case 24:
		return Item{Code: "GU1", Name: "Mark-1 Gun Unit", Quantity: qty}
	case 25:
		return Item{Code: "GU2", Name: "Mark-2 Gun Unit", Quantity: qty}
	case 26:
		return Item{Code: "GU3", Name: "Mark-3 Gun Unit", Quantity: qty}
	case 27:
		return Item{Code: "GU4", Name: "Mark-4 Gun Unit", Quantity: qty}
	case 28:
		return Item{Code: "GU5", Name: "Mark-5 Gun Unit", Quantity: qty}
	case 29:
		return Item{Code: "GU6", Name: "Mark-6 Gun Unit", Quantity: qty}
	case 30:
		return Item{Code: "GU7", Name: "Mark-7 Gun Unit", Quantity: qty}
	case 31:
		return Item{Code: "GU8", Name: "Mark-8 Gun Unit", Quantity: qty}
	case 32:
		return Item{Code: "GU9", Name: "Mark-9 Gun Unit", Quantity: qty}
	case 33:
		return Item{Code: "X1", Name: "X1 Unit", Quantity: qty}
	case 34:
		return Item{Code: "X2", Name: "X2 Unit", Quantity: qty}
	case 35:
		return Item{Code: "X3", Name: "X3 Unit", Quantity: qty}
	case 36:
		return Item{Code: "X4", Name: "X4 Unit", Quantity: qty}
	case 37:
		return Item{Code: "X5", Name: "X5 Unit", Quantity: qty}
	default:
		return Item{}
	}
}

// codeToCargoCapacity returns cargo capacity based on class and tonnage
func codeToShipCargoCapacity(code int, tonnage int) int {
	switch code {
	case 0: // PB
		return 1
	case 1: // CT
		return 2
	case 2: // ES
		return 5
	case 3: // FF
		return 10
	case 4: // DD
		return 15
	case 5: // CL
		return 20
	case 6: // CS
		return 25
	case 7: // CA
		return 30
	case 8: // CC
		return 35
	case 9: // BC
		return 40
	case 10: // BS
		return 45
	case 11: // DN
		return 50
	case 12: // SD
		return 55
	case 13: // BM
		return 60
	case 14: // BW
		return 65
	case 15: // BS
		return 70
	case 16: // Starbase
		return tonnage
	case 17: // Transport
		return (10 + (tonnage / 2)) * tonnage
	default:
		return 0
	}
}

func codeToShipClass(code int) string {
	switch code {
	case 0:
		return "PB"
	case 1:
		return "CT"
	case 2:
		return "ES"
	case 3:
		return "FF"
	case 4:
		return "DD"
	case 5:
		return "CL"
	case 6:
		return "CS"
	case 7:
		return "CA"
	case 8:
		return "CC"
	case 9:
		return "BC"
	case 10:
		return "BS"
	case 11:
		return "DN"
	case 12:
		return "SD"
	case 13:
		return "BM"
	case 14:
		return "BW"
	case 15:
		return "BR"
	case 16:
		return "BA"
	case 17:
		return "TR"
	default:
		return ""
	}
}

// codeToShipCost returns the original cost to build the ship
func codeToShipCost(code int, tonnage int) int {
	switch code {
	case 0: // PB
		return 100
	case 1: // CT
		return 200
	case 2: // ES
		return 500
	case 3: // FF
		return 1_000
	case 4: // DD
		return 1_500
	case 5: // CL
		return 2_000
	case 6: // CS
		return 2_500
	case 7: // CA
		return 3_000
	case 8: // CC
		return 3_500
	case 9: // BC
		return 4_000
	case 10: // BS
		return 4_500
	case 11: // DN
		return 5_000
	case 12: // SD
		return 5_500
	case 13: // BM
		return 6_000
	case 14: // BW
		return 6_500
	case 15: // BS
		return 7_000
	case 16: // Starbase
		return 100 * tonnage
	case 17: // Transport
		return 100 * tonnage
	default:
		return 0
	}
}

// codeToShipSize returns the ship size based on class and tonnage
func codeToShipSize(code int, tonnage int) int {
	switch code {
	case 17:
		return tonnage
	default:
		return 0
	}
}

// codeToShipTonnage returns the ship size in tons based on class and tonnage
func codeToShipTonnage(code int, tonnage int) int {
	switch code {
	case 0: // PB
		return 10_000
	case 1: // CT
		return 20_000
	case 2: // ES
		return 50_000
	case 3: // FF
		return 100_000
	case 4: // DD
		return 150_000
	case 5: // CL
		return 200_000
	case 6: // CS
		return 250_000
	case 7: // CA
		return 300_000
	case 8: // CC
		return 350_000
	case 9: // BC
		return 400_000
	case 10: // BS
		return 450_000
	case 11: // DN
		return 500_000
	case 12: // SD
		return 550_000
	case 13: // BM
		return 600_000
	case 14: // BW
		return 650_000
	case 15: // BS
		return 700_000
	case 16: // Starbase
		return 10_000 * tonnage
	case 17: // Transport
		return 10_000 * tonnage
	default:
		return 0
	}
}

func codeToStarColor(code int) StarColor {
	switch code {
	case 1:
		return StarColor{Code: "O", Name: "Blue"}
	case 2:
		return StarColor{Code: "B", Name: "Blue-White"}
	case 3:
		return StarColor{Code: "A", Name: "White"}
	case 4:
		return StarColor{Code: "F", Name: "Yellow-White"}
	case 5:
		return StarColor{Code: "G", Name: "Yellow"}
	case 6:
		return StarColor{Code: "K", Name: "Orange"}
	case 7:
		return StarColor{Code: "M", Name: "Red"}
	default:
		return StarColor{}
	}
}

func codeToStarType(code int) StarType {
	switch code {
	case 1:
		return StarType{Code: "d", Name: "Dwarf"}
	case 2:
		return StarType{Code: "D", Name: "Degenerate"}
	case 3:
		return StarType{Code: " ", Name: "Main Sequence"}
	case 4:
		return StarType{Code: "G", Name: "Giant"}
	default:
		return StarType{}
	}
}

func nameToString(name [32]uint8) string {
	var b []byte
	for _, ch := range name {
		if ch == 0 {
			break
		}
		b = append(b, byte(ch))
	}
	return string(b)
}

// speciesBitIsSet returns true if the bit is set for the species.
// note: the species number must be 1 based!
// sp01       65536                       1 0000 0000 0000 0000
// sp09    16777216             1 0000 0000 0000 0000 0000 0000
// sp18  8589934592  10 0000 0000 0000 0000 0000 0000 0000 0000
func speciesBitIsSet(set [2]uint64, sp int) bool {
	return (set[0] & (1 << (sp + 15))) != 0
}
