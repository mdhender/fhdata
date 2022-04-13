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
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

// readGalaxy returns either an initialized galaxy_data or an error.
func readGalaxy(name string, bo binary.ByteOrder) (*galaxy_data, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	g := &galaxy_data{}
	if err := binary.Read(r, bo, g); err != nil {
		return nil, err
	}

	return g, nil
}

func readNamplas(r io.Reader, num_namplas int, bo binary.ByteOrder) ([]nampla_data, error) {
	namplas := make([]nampla_data, num_namplas)
	for i := 0; i < num_namplas; i++ {
		if err := binary.Read(r, bo, &namplas[i]); err != nil {
			return nil, err
		}
	}
	return namplas, nil
}

// readPlanets returns either an initialized set of planets or an error
func readPlanets(name string, bo binary.ByteOrder) ([]planet_data, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var pd planet_file
	if err := binary.Read(r, bo, &pd.NumPlanets); err != nil {
		return nil, err
	}

	num_planets := int(pd.NumPlanets)
	planet_base := make([]planet_data, num_planets, num_planets)
	for i := 0; i < num_planets; i++ {
		if err := binary.Read(r, bo, &planet_base[i]); err != nil {
			return nil, err
		}
	}

	return planet_base, nil
}

func readShips(r io.Reader, num_ships int, bo binary.ByteOrder) ([]ship_data, error) {
	ships := make([]ship_data, num_ships)
	for i := 0; i < num_ships; i++ {
		if err := binary.Read(r, bo, &ships[i]); err != nil {
			return nil, err
		}
	}
	return ships, nil
}

// readSpecies returns either an initialized species with namplas and ships or an error.
func readSpecies(name string, bo binary.ByteOrder) (*species_file, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	sp := &species_file{
		data: &species_data{},
	}
	if err := binary.Read(r, bo, sp.data); err != nil {
		return nil, err
	}

	sp.namplas, err = readNamplas(r, int(sp.data.NumNamplas), bo)
	if err != nil {
		return nil, err
	}

	sp.ships, err = readShips(r, int(sp.data.NumShips), bo)
	if err != nil {
		return nil, err
	}

	return sp, nil
}

// readStars returns either an initialized set of star_data or an error.
func readStars(name string, bo binary.ByteOrder) ([]star_data, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var sd star_file
	if err := binary.Read(r, bo, &sd.NumStars); err != nil {
		return nil, err
	}

	num_stars := int(sd.NumStars)
	star_base := make([]star_data, num_stars, num_stars)
	for i := 0; i < num_stars; i++ {
		if err := binary.Read(r, bo, &star_base[i]); err != nil {
			return nil, err
		}
	}

	return star_base, nil
}
