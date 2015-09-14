/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@opensource.org>, 2015
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>. }}} */

package license

import (
	"encoding/json"
	"fmt"
	"os"
)

//
type Identifier struct {
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme"`
}

//
type Link struct {
	Note *string `json:"note"`
	URL  string  `json:"url"`
}

//
type OtherName struct {
	Name string  `json:"name"`
	Note *string `json:"note"`
}

//
type Text struct {
	ContentType string `json:"content_type"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

//
type License struct {
	Id           string       `json:"id"`
	Identifiers  []Identifier `json:"identifiers"`
	Links        []Link       `json:"links"`
	Name         string       `json:"name"`
	OtherNames   []OtherName  `json:"other_names"`
	SupersededBy *string      `json:"superseded_by"`
	Tags         []string     `json:"tags"`
	Texts        []Text       `json:"text"`
}

//
type Licenses []License

//
func (licenses Licenses) GetIds() []string {
	identifiers := []string{}
	for _, license := range licenses {
		identifiers = append(identifiers, license.Id)
	}
	return identifiers
}

//
func (licenses Licenses) GetIdMap() map[string]License {
	ret := map[string]License{}
	for _, license := range licenses {
		ret[license.Id] = license
		for _, identifier := range license.Identifiers {
			ret[fmt.Sprintf("%s/%s", identifier.Scheme, identifier.Identifier)] = license
		}
	}
	return ret
}

//
func (licenses Licenses) GetTagMap() map[string][]License {
	ret := map[string][]License{}
	for _, license := range licenses {
		for _, tag := range license.Tags {
			ret[tag] = append(ret[tag], license)
		}
	}
	return ret
}

//
func LoadLicensesFiles(path string) (Licenses, error) {
	ret := Licenses{}
	fh, err := os.Open(path)
	if err != nil {
		return Licenses{}, err
	}
	if err := json.NewDecoder(fh).Decode(&ret); err != nil {
		return Licenses{}, err
	}
	return ret, nil
}

// vim: foldmethod=marker
