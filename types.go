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

import "fmt"

type AtmosphericGas struct {
	Gas
	Pct int
}

type Cluster struct {
	Turn               int
	Radius             int
	DesignedNumSpecies int
	Planets            []*Planet
	Species            []*Species
	Systems            []*System
}

type Coords struct {
	X int
	Y int
	Z int
}

func (c Coords) Equals(o Coords) bool {
	return c.X == o.X && c.Y == o.Y && c.Z == o.Z
}

func (c Coords) Less(o Coords) bool {
	if c.X < o.X {
		return true
	} else if c.X > o.X {
		return false
	}
	if c.Y < o.Y {
		return true
	} else if c.Y > o.Y {
		return false
	}
	return c.Z < o.Z
}

func (c Coords) String() string {
	return fmt.Sprintf("%d, %d, %d", c.X, c.Y, c.Z)
}

type Colony struct {
	Id         int
	Coords     Coords
	DevelopAUs *Develop
	DevelopIUs *Develop
	Inventory  []Item
	Is         struct {
		Colony          bool
		DisbandedColony bool
		Hidden          bool
		Hiding          bool
		HomePlanet      bool
		HomeWorld       bool
		MiningColony    bool
		Populated       bool
		ResortColony    bool
	}
	LSN               int
	ManufacturingBase int
	Message           int
	MiningBase        int
	Name              string
	Orbit             int // the first orbit is 1
	Planet            *Planet
	PopulationUnits   int
	Production        int
	Shipyards         int
	SiegeEffPct       int
	Species           *Species
	Special           int
	Status            int
	System            *System
	UseOnAmbush       int
}

type Develop struct {
	Code           string
	AutoInstall    int
	UnitsNeeded    int
	UnitsToInstall int
}

type Gas struct {
	Code string
	Name string
}

type Item struct {
	Code     string
	Name     string
	Quantity int
}

type Location struct {
	Colony *Colony
	Planet *Planet
	System *System
}

type Planet struct {
	Id             int
	Atmosphere     []*AtmosphericGas
	Colonies       []*Colony
	Coords         Coords
	Diameter       int
	EconEfficiency int
	Gravity        int
	Is             struct {
		IdealColonyPlanet   bool
		IdealHomePlanet     bool
		RadioactiveHellHole bool
	}
	Message                  int
	MiningDifficultyBase     int
	MiningDifficultyIncrease int
	Orbit                    int // the first orbit is 1
	PressureClass            int
	System                   *System
	TemperatureClass         int
}

type Ship struct {
	Id                 int
	Age                int
	ArrivedViaWormhole bool
	Class              string
	CargoCapacity      int
	Coords             Coords
	Destination        *Location
	ForcedJump         bool
	Hiding             bool
	InDeepSpace        bool
	InOrbit            bool
	Inventory          []Item
	JumpedInCombat     bool
	JustJumped         bool
	LoadingPoint       *Colony
	Location           Location
	Name               string
	OnSurface          bool
	Orbit              int // the first orbit is 1
	RemainingCost      int
	Size               int // meaningful only for transports
	Special            int
	Species            *Species
	Status             int
	SubLight           bool
	Tonnage            int
	TotalCost          int
	UnderConstruction  bool
	UnloadingPoint     *Colony
}

type Species struct {
	Id                   int
	Allies               map[string]*Species
	AutoOrders           bool
	Colonies             []*Colony
	Contacts             map[string]*Species
	EconUnitsBanked      int
	EconUnitsProduced    int
	Enemies              map[string]*Species
	FleetMaintenanceCost int
	FleetMaintenancePct  int
	Gases                struct {
		Neutral  []Gas
		Poison   []Gas
		Required struct {
			Gas
			MinPct int
			MaxPct int
		}
	}
	GovtName               string
	GovtType               string
	HomeColony             *Colony
	HomePlanet             *Planet
	HomePlanetOriginalBase int
	HomeSystem             *System
	Name                   string
	Ships                  []*Ship
	SystemsScanned         []*System
	SystemsVisited         []*System
	BI, GV, LS, MA, MI, ML Tech
}

type StarColor struct {
	Code string
	Name string
}

type StarType struct {
	Code string
	Name string
}

type System struct {
	Id     int
	Color  StarColor
	Coords Coords
	Is     struct {
		HomeSystem bool
	}
	Message      int
	Planets      []*Planet
	ScannedBy    map[string]*Species
	Size         int
	Type         StarType
	VisitedBy    map[string]*Species
	WormholeExit *System
}

type Tech struct {
	Code           string
	CurrentLevel   int
	InitialLevel   int
	KnowledgeLevel int
	Name           string
	XPs            int
}

const (
	MAX_ITEMS = 38
	// Status code of named planet. These are logically ORed together.
	HOME_PLANET      = 1
	COLONY           = 2
	POPULATED        = 8
	MINING_COLONY    = 16
	RESORT_COLONY    = 32
	DISBANDED_COLONY = 64
	// Ship status codes.
	UNDER_CONSTRUCTION = 0
	ON_SURFACE         = 1
	IN_ORBIT           = 2
	IN_DEEP_SPACE      = 3
	JUMPED_IN_COMBAT   = 4
	FORCED_JUMP        = 5
)

// galaxy_data is the layout in the binary data file.
type galaxy_data struct {
	/* Design number of species in galaxy. */
	DNumSpecies int32
	/* Actual number of species allocated. */
	NumSpecies int32
	/* Galactic radius in parsecs. */
	Radius int32
	/* Current turn number. */
	TurnNumber int32
}

// nampla_data is the layout in the binary data file.
type nampla_data struct {
	/* Name of planet. */
	Name [32]uint8
	/* Coordinates. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Status code of named planet. These are logically ORed together. */
	Status uint8
	/* Reserved for future use. Zero for now. */
	Reserved1 uint8
	/* HIDE order given. */
	Hiding uint8
	/* Colony is hidden. */
	Hidden uint8
	/* Reserved for future use. Zero for now. */
	Reserved2 int16
	/* Index (starting at zero) into the file "planets.dat" of this planet. */
	PlanetIndex int16
	/* Siege effectiveness - a percentage between 0 and 99. */
	SiegeEff int16
	/* Number of shipyards on planet. */
	Shipyards int16
	/* Reserved for future use. Zero for now. */
	Reserved4 int32
	/* Incoming ship with only CUs on board. */
	IUsNeeded int32
	/* Incoming ship with only CUs on board. */
	AUsNeeded int32
	/* Number of IUs to be automatically installed. */
	AutoIUs int32
	/* Number of AUs to be automatically installed. */
	AutoAUs int32
	/* Reserved for future use. Zero for now. */
	Reserved5 int32
	/* Colonial mining units to be installed. */
	IUsToInstall int32
	/* Colonial manufacturing units to be installed. */
	AUsToInstall int32
	/* Mining base times 10. */
	MiBase int32
	/* Manufacturing base times 10. */
	MaBase int32
	/* Number of available population units. */
	PopUnits int32
	/* Quantity of each item available. */
	ItemQuantity [MAX_ITEMS]int32
	/* Reserved for future use. Zero for now. */
	Reserved6 int32
	/* Amount to use on ambush. */
	UseOnAmbush int32
	/* Message associated with this planet, if any. */
	Message int32
	/* Different for each application. */
	Special int32
	/* Use for expansion. Initialized to all zeroes. */
	Padding [28]uint8
}

// planet_data is the layout in the binary data file.
type planet_data struct {
	/* Temperature class, 1-30. */
	TemperatureClass int8
	/* Pressure class, 0-29. */
	PressureClass int8
	/* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Special int8
	/* Reserved for future use. Zero for now. */
	Reserved1 int8
	/* Gas in atmosphere. Zero if none. */
	Gas [4]int8
	/* Percentage of gas in atmosphere. */
	GasPercent [4]int8
	/* Reserved for future use. Zero for now. */
	Reserved2 int16
	/* Diameter in thousands of kilometers. */
	Diameter int16
	/* Surface gravity. Multiple of Earth gravity times 100. */
	Gravity int16
	/* Mining difficulty times 100. */
	MiningDifficulty int16
	/* Economic efficiency. Always 100 for a home planet. */
	EconEfficiency int16
	/* Increase in mining difficulty. */
	MDIncrease int16
	/* Message associated with this planet, if any. */
	Message int32
	/* Reserved for future use. Zero for now. */
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
}

// planet_file is a helper struct that represents the layout of data in the binary file
type planet_file struct {
	NumPlanets int32
	PlanetBase []planet_data
}

// ship_data is the layout in the binary data file.
type ship_data struct {
	/* Name of ship. */
	Name [32]uint8
	/* Current coordinates. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Current status of ship. */
	Status uint8
	/* Ship type. */
	Type uint8
	/* Destination if ship was forced to jump from combat. Also used by TELESCOPE command. */
	DestX uint8
	DestY uint8
	DestZ uint8
	/* Set if ship jumped this turn. */
	JustJumped uint8
	/* Ship arrived via wormhole in the PREVIOUS turn. */
	ArrivedViaWormhole uint8
	/* Reserved for future use. Zero for now. */
	Reserved1 uint8
	Reserved2 int16
	Reserved3 int16
	/* Ship class. */
	Class int16
	/* Ship tonnage divided by 10,000. */
	Tonnage int16
	/* Quantity of each item carried. */
	ItemQuantity [MAX_ITEMS]int16
	/* Ship age. */
	Age int16
	/* The cost needed to complete the ship if still under construction. */
	RemainingCost int16
	/* Reserved for future use. Zero for now. */
	Reserved4 int16
	/* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	LoadingPoint int16
	/* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	UnloadingPoint int16
	/* Different for each application. */
	Special int32
	/* Use for expansion. Initialized to all zeroes. */
	Padding [28]uint8
	// padding to make Go struct same size as C
	MorePadding [2]uint8
}

// species_data is the layout in the binary data file.
type species_data struct {
	/* Name of species. */
	Name [32]uint8
	/* Name of government. */
	GovtName [32]uint8
	/* Type of government. */
	GovtType [32]uint8
	/* Coordinates of home planet. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Gas required by species. */
	RequiredGas uint8
	/* Minimum needed percentage. */
	RequiredGasMin uint8
	/* Maximum allowed percentage. */
	RequiredGasMax uint8
	/* Reserved for future use. Zero for now. */
	Reserved5 uint8
	/* Gases neutral to species. */
	NeutralGas [6]uint8
	/* Gases poisonous to species. */
	PoisonGas [6]uint8
	/* AUTO command was issued. */
	AutoOrders uint8
	/* Reserved for future use. Zero for now. */
	Reserved3 uint8
	Reserved4 int16
	/* Actual tech levels. */
	TechLevel [6]int16
	/* Tech levels at start of turn. */
	InitTechLevel [6]int16
	/* Unapplied tech level knowledge. */
	TechKnowledge [6]int16
	/* Number of named planets, including home planet and colonies. */
	NumNamplas int32
	/* Number of ships. */
	NumShips int32
	/* Experience points for tech levels. */
	TechEps [6]int32
	/* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	HPOriginalBase int32
	/* Number of economic units. */
	EconUnits int32
	/* Total fleet maintenance cost. */
	FleetCost int32
	/* Fleet maintenance cost as a percentage times one hundred. */
	FleetPercentCost int32
	/* A bit is set if corresponding species has been met. */
	Contact [2]uint64
	/* A bit is set if corresponding species is considered an ally. */
	Ally [2]uint64
	/* A bit is set if corresponding species is considered an enemy. */
	Enemy [2]uint64
	/* Use for expansion. Initialized to all zeroes. */
	Padding [12]uint8
}

// species_file is a helper struct for loading the combined species data
type species_file struct {
	data    *species_data
	namplas []nampla_data
	ships   []ship_data
}

// star_data is the layout in the binary data file.
type star_data struct {
	/* Coordinates. */
	X int8
	Y int8
	Z int8
	/* Dwarf, degenerate, main sequence or giant. */
	Type int8
	/* Star color. Blue, blue-white, etc. */
	Color int8
	/* Star size, from 0 thru 9 inclusive. */
	Size int8
	/* Number of usable planets in star system. */
	NumPlanets int8
	/* TRUE if this is a good potential home system. */
	HomeSystem int8
	/* TRUE if wormhole entry/exit. */
	WormHere int8
	/* Coordinates of exit point for wormhole. Valid only if WormHere is TRUE. */
	WormX int8
	WormY int8
	WormZ int8
	/* Reserved for future use. Zero for now. */
	Reserved1 int16
	Reserved2 int16
	/* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	PlanetIndex int16
	/* Message associated with this star system, if any. */
	Message int32
	/* A bit is set if corresponding species has been here. */
	VisitedBy [2]uint64
	/* Reserved for future use. Zero for now. */
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
	// padding to make Go struct same size as C
	Padding [2]uint8
}

// star_file is a helper struct that represents the layout in the binary data file.
type star_file struct {
	NumStars int32
	StarBase []star_data
}
