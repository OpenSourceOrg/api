/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@opensource.org>, 2015-2016
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

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
	MediaType string `json:"media_type"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}

//
type License struct {
	Id           string       `json:"id"`
	Identifiers  []Identifier `json:"identifiers"`
	Links        []Link       `json:"links"`
	Name         string       `json:"name"`
	OtherNames   []OtherName  `json:"other_names"`
	SupersededBy *string      `json:"superseded_by"`
	Keywords     []string     `json:"keywords"`
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
		for _, tag := range license.Keywords {
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
