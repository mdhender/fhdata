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

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
)

// LoadFromPath loads the galaxy, stars, planets, and species files from the given path.
func LoadFromPath(dataPath string, bo binary.ByteOrder) (*Cluster, error) {
	// load the galaxy, stars, planets, and species data from the binary files.
	galaxy, err := readGalaxy(filepath.Join(dataPath, "galaxy.dat"), bo)
	if err != nil {
		return nil, err
	}
	stars, err := readStars(filepath.Join(dataPath, "stars.dat"), bo)
	if err != nil {
		return nil, err
	}
	planets, err := readPlanets(filepath.Join(dataPath, "planets.dat"), bo)
	if err != nil {
		return nil, err
	}
	var speciesData []*species_file
	for i := 0; i < int(galaxy.NumSpecies); i++ {
		spNo := i + 1
		sp, err := readSpecies(filepath.Join(dataPath, fmt.Sprintf("sp%02d.dat", spNo)), bo)
		if err != nil {
			return nil, err
		}
		speciesData = append(speciesData, sp)
	}

	// translate the galaxy data
	cluster := &Cluster{}
	cluster.DesignedNumSpecies = int(galaxy.DNumSpecies)
	cluster.Planets = make([]*Planet, len(planets), len(planets))
	cluster.Radius = int(galaxy.Radius)
	cluster.Species = make([]*Species, int(galaxy.NumSpecies), int(galaxy.NumSpecies))
	cluster.Systems = make([]*System, len(stars), len(stars))
	cluster.Turn = int(galaxy.TurnNumber)

	// translate the star data
	for i, star := range stars {
		system := &System{
			Id:        i + 1,
			Color:     codeToStarColor(int(star.Color)),
			Coords:    Coords{X: int(star.X), Y: int(star.Y), Z: int(star.Z)},
			Message:   int(star.Message),
			Planets:   make([]*Planet, int(star.NumPlanets), int(star.NumPlanets)),
			ScannedBy: make(map[string]*Species),
			Size:      int(star.Size),
			Type:      codeToStarType(int(star.Type)),
			VisitedBy: make(map[string]*Species),
		}
		system.Is.HomeSystem = star.HomeSystem != 0

		// create planets in systems and add coordinates and orbit values
		for pn := 0; pn < len(system.Planets); pn++ {
			planet := &Planet{
				Id:     int(star.PlanetIndex) + pn + 1,
				Coords: system.Coords,
				Orbit:  pn + 1,
				System: system,
			}
			system.Planets[pn] = planet
			cluster.Planets[int(star.PlanetIndex)+pn] = planet
		}

		cluster.Systems[i] = system
	}

	// link wormholes
	for i, star := range stars {
		if star.WormHere == 0 {
			continue
		}
		coords := Coords{X: int(star.WormX), Y: int(star.WormY), Z: int(star.WormZ)}
		for _, system := range cluster.Systems {
			if coords.Equals(system.Coords) {
				cluster.Systems[i].WormholeExit = system
				system.WormholeExit = cluster.Systems[i]
				break
			}
		}
	}

	// translate the planet data
	for pn, pp := range planets {
		planet := cluster.Planets[pn]
		planet.Diameter = int(pp.Diameter)
		planet.EconEfficiency = int(pp.EconEfficiency)
		planet.Gravity = int(pp.Gravity)
		planet.Message = int(pp.Message)
		planet.MiningDifficultyBase = int(pp.MiningDifficulty)
		planet.MiningDifficultyIncrease = int(pp.MDIncrease)
		planet.PressureClass = int(pp.PressureClass)
		planet.TemperatureClass = int(pp.TemperatureClass)
		for i, code := range pp.Gas {
			if pp.GasPercent[i] != 0 {
				planet.Atmosphere = append(planet.Atmosphere, &AtmosphericGas{Gas: codeToGas(int(code)), Pct: int(pp.GasPercent[i])})
			}
		}
		switch int(pp.Special) {
		case 1:
			planet.Is.IdealHomePlanet = true
		case 2:
			planet.Is.IdealColonyPlanet = true
		case 3:
			planet.Is.RadioactiveHellHole = true
		}
	}

	// translate the species data
	for i, sp := range speciesData {
		spNo := i + 1
		species := sp.data
		cluster.Species[i] = &Species{
			Id:                     spNo,
			Allies:                 make(map[string]*Species),
			AutoOrders:             species.AutoOrders != 0,
			Contacts:               make(map[string]*Species),
			EconUnitsBanked:        int(species.EconUnits),
			Enemies:                make(map[string]*Species),
			FleetMaintenanceCost:   int(species.FleetCost),
			FleetMaintenancePct:    int(species.FleetPercentCost),
			GovtName:               nameToString(species.GovtName),
			GovtType:               nameToString(species.GovtType),
			HomePlanetOriginalBase: int(species.HPOriginalBase),
			Name:                   nameToString(species.Name),
		}
		coords, orbit := Coords{X: int(species.X), Y: int(species.Y), Z: int(species.Z)}, int(species.PN)
		for _, planet := range cluster.Planets {
			if coords.Equals(planet.Coords) && planet.Orbit == orbit {
				cluster.Species[i].HomePlanet = planet
				cluster.Species[i].HomeSystem = planet.System
				break
			}
		}
		for _, code := range species.NeutralGas {
			if code != 0 {
				cluster.Species[i].Gases.Neutral = append(cluster.Species[i].Gases.Neutral, codeToGas(int(code)))
			}
		}
		for _, code := range species.PoisonGas {
			if code != 0 {
				cluster.Species[i].Gases.Poison = append(cluster.Species[i].Gases.Poison, codeToGas(int(code)))
			}
		}
		cluster.Species[i].Gases.Required.Gas = codeToGas(int(species.RequiredGas))
		cluster.Species[i].Gases.Required.MinPct = int(species.RequiredGasMin)
		cluster.Species[i].Gases.Required.MaxPct = int(species.RequiredGasMax)
		cluster.Species[i].BI = Tech{
			Code:           "BI",
			Name:           "Biology",
			CurrentLevel:   int(species.TechLevel[5]),
			InitialLevel:   int(species.InitTechLevel[5]),
			KnowledgeLevel: int(species.TechKnowledge[5]),
			XPs:            int(species.TechEps[5]),
		}
		cluster.Species[i].GV = Tech{
			Code:           "GV",
			Name:           "Gravitics",
			CurrentLevel:   int(species.TechLevel[3]),
			InitialLevel:   int(species.InitTechLevel[3]),
			KnowledgeLevel: int(species.TechKnowledge[3]),
			XPs:            int(species.TechEps[3]),
		}
		cluster.Species[i].LS = Tech{
			Code:           "LS",
			Name:           "Life Support",
			CurrentLevel:   int(species.TechLevel[4]),
			InitialLevel:   int(species.InitTechLevel[4]),
			KnowledgeLevel: int(species.TechKnowledge[4]),
			XPs:            int(species.TechEps[4]),
		}
		cluster.Species[i].MA = Tech{
			Code:           "MA",
			Name:           "Manufacturing",
			CurrentLevel:   int(species.TechLevel[1]),
			InitialLevel:   int(species.InitTechLevel[1]),
			KnowledgeLevel: int(species.TechKnowledge[1]),
			XPs:            int(species.TechEps[1]),
		}
		cluster.Species[i].MI = Tech{
			Code:           "MI",
			Name:           "Mining",
			CurrentLevel:   int(species.TechLevel[0]),
			InitialLevel:   int(species.InitTechLevel[0]),
			KnowledgeLevel: int(species.TechKnowledge[0]),
			XPs:            int(species.TechEps[0]),
		}
		cluster.Species[i].ML = Tech{
			Code:           "ML",
			Name:           "Military",
			CurrentLevel:   int(species.TechLevel[2]),
			InitialLevel:   int(species.InitTechLevel[2]),
			KnowledgeLevel: int(species.TechKnowledge[2]),
			XPs:            int(species.TechEps[2]),
		}

		// translate the species colony data
		for n, nampla := range sp.namplas {
			colony := &Colony{
				Id:                n + 1,
				Coords:            Coords{X: int(nampla.X), Y: int(nampla.Y), Z: int(nampla.Z)},
				ManufacturingBase: int(nampla.MaBase),
				Message:           int(nampla.Message),
				MiningBase:        int(nampla.MiBase),
				Name:              nameToString(nampla.Name),
				Orbit:             int(nampla.PN),
				PopulationUnits:   int(nampla.PopUnits),
				Shipyards:         int(nampla.Shipyards),
				SiegeEffPct:       int(nampla.SiegeEff),
				Species:           cluster.Species[i],
				UseOnAmbush:       int(nampla.UseOnAmbush),
			}
			cluster.Species[i].Colonies = append(cluster.Species[i].Colonies, colony)

			if nampla.AutoAUs != 0 || nampla.AUsNeeded != 0 || nampla.AUsToInstall != 0 {
				colony.DevelopAUs = &Develop{
					AutoInstall:    int(nampla.AutoAUs),
					UnitsNeeded:    int(nampla.AUsNeeded),
					UnitsToInstall: int(nampla.AUsToInstall),
				}
			}
			if nampla.AutoIUs != 0 || nampla.IUsNeeded != 0 || nampla.IUsToInstall != 0 {
				colony.DevelopIUs = &Develop{
					AutoInstall:    int(nampla.AutoIUs),
					UnitsNeeded:    int(nampla.IUsNeeded),
					UnitsToInstall: int(nampla.IUsToInstall),
				}
			}
			for code, qty := range nampla.ItemQuantity {
				if qty != 0 {
					colony.Inventory = append(colony.Inventory, codeToItem(code, int(qty)))
				}
			}
			colony.Is.Colony = (nampla.Status & COLONY) != 0
			colony.Is.DisbandedColony = (nampla.Status & DISBANDED_COLONY) != 0
			colony.Is.Hidden = nampla.Hidden != 0
			colony.Is.Hiding = nampla.Hiding != 0
			colony.Is.HomePlanet = (nampla.Status & HOME_PLANET) != 0
			colony.Is.HomeWorld = n == 0
			colony.Is.MiningColony = (nampla.Status & MINING_COLONY) != 0
			colony.Is.Populated = (nampla.Status & POPULATED) != 0
			colony.Is.ResortColony = (nampla.Status & RESORT_COLONY) != 0

			// link the colony to the planet and system it occupies
			for _, planet := range cluster.Planets {
				if colony.Coords.Equals(planet.Coords) && colony.Orbit == planet.Orbit {
					colony.Planet = planet
					colony.System = planet.System
					planet.Colonies = append(planet.Colonies, colony)
					break
				}
			}

			if colony.Is.HomeWorld {
				// link the species home world
				cluster.Species[i].HomeColony = colony
			}
		}

		// translate the species ship data
		for n, sh := range sp.ships {
			ship := &Ship{
				Id:                 n + 1,
				Age:                int(sh.Age),
				ArrivedViaWormhole: sh.ArrivedViaWormhole != 0,
				CargoCapacity:      codeToShipCargoCapacity(int(sh.Class), int(sh.Tonnage)),
				Class:              codeToShipClass(int(sh.Class)),
				Coords:             Coords{X: int(sh.X), Y: int(sh.Y), Z: int(sh.Z)},
				ForcedJump:         sh.Status == FORCED_JUMP,
				InDeepSpace:        sh.Status == IN_DEEP_SPACE,
				InOrbit:            sh.Status == IN_ORBIT,
				JumpedInCombat:     sh.Status == JUMPED_IN_COMBAT,
				JustJumped:         sh.JustJumped != 0,
				Name:               nameToString(sh.Name),
				OnSurface:          sh.Status == ON_SURFACE,
				Orbit:              int(sh.PN),
				RemainingCost:      int(sh.RemainingCost),
				Size:               codeToShipSize(int(sh.Class), int(sh.Tonnage)),
				Species:            cluster.Species[i],
				SubLight:           sh.Type != 0 || sh.Class == 16, /* sublight or starbase */
				Tonnage:            codeToShipTonnage(int(sh.Class), int(sh.Tonnage)),
				UnderConstruction:  sh.Status == UNDER_CONSTRUCTION,
			}
			cluster.Species[i].Ships = append(cluster.Species[i].Ships, ship)

			for code, qty := range sh.ItemQuantity {
				if qty != 0 {
					ship.Inventory = append(ship.Inventory, codeToItem(code, int(qty)))
				}
			}
			for _, system := range cluster.Systems {
				if ship.Coords.Equals(system.Coords) {
					ship.Location.System = system
					if ship.Orbit != 0 {
						ship.Location.Planet = system.Planets[ship.Orbit-1]
					}
					break
				}
			}
		}
	}

	// link species to systems they've visited or scanned.
	// we assume that every system that has a colony or ship in it is being scanned.
	for i, star := range stars {
		system := cluster.Systems[i]
		for n, species := range cluster.Species {
			spNo := n + 1
			if hasVisited := speciesBitIsSet(star.VisitedBy, spNo); hasVisited {
				system.VisitedBy[species.Name] = species
				species.SystemsVisited = append(species.SystemsVisited, system)
			}
			hasScanned := false
			for _, colony := range species.Colonies {
				if colony.System == system {
					hasScanned = true
					break
				}
			}
			if !hasScanned {
				for _, ship := range species.Ships {
					if ship.Location.System == system {
						hasScanned = true
						break
					}
				}
			}
			if hasScanned {
				system.ScannedBy[species.Name] = species
				species.SystemsScanned = append(species.SystemsScanned, system)
			}
		}
	}

	// calculate the amount of life support needed for each species and planet
	for _, planet := range cluster.Planets {
		planet.LSN = make([]int, len(cluster.Species), len(cluster.Species))
		for i, species := range cluster.Species {
			// assuming required gas is NOT present and so requires 3 points of life support
			planet.LSN[i] = 3
			// temperature class requires 3 points of LS per point of difference
			tc := planet.TemperatureClass - species.HomePlanet.TemperatureClass
			if tc < 0 {
				tc = -tc
			}
			// each point of difference requires 3 points of life support
			planet.LSN[i] += 3 * tc
			// pressure class requires 3 points of LS per point of difference
			pc := planet.PressureClass - species.HomePlanet.PressureClass
			if pc < 0 {
				pc = -pc
			}
			planet.LSN[i] += 3 * pc
			for _, atmo := range planet.Atmosphere {
				// check if the slot has gas
				if atmo.Pct != 0 {
					// check if required gas is present
					if atmo.Code == species.Gases.Required.Code {
						// and in the right amount
						if species.Gases.Required.MinPct <= atmo.Pct && atmo.Pct <= species.Gases.Required.MaxPct {
							// required gas is present and in the right amount, so undo 3 points of life support
							planet.LSN[i] -= 3
						}
					} else {
						// check if it is a poisonous gas
						for _, poison := range species.Gases.Poison {
							if atmo.Code == poison.Code {
								// each poisonous gas requires 3 points of life support
								planet.LSN[i] += 3
								break
							}
						}
					}
				}
			}
		}
	}

	// copy the life support into the colonies
	for i, species := range cluster.Species {
		for _, colony := range species.Colonies {
			colony.LSN = colony.Planet.LSN[i]
		}
	}
	return cluster, nil
}
