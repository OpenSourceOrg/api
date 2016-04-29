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

package client

import (
	"github.com/opensourceorg/api/license"

	"encoding/json"
	"net/http"
	"strings"
)

var baseURL = "https://api.opensource.org"

func createURIString(resource string) string {
	/* XXX: TODO use URL path stuff */
	return baseURL + "/" + resource
}

type BadRequest struct {
	errors []map[string]string
}

func (b BadRequest) messages() []string {
	ret := []string{}
	for _, err := range b.errors {
		ret = append(ret, err["message"])
	}
	return ret
}

func (b BadRequest) Error() string {
	return strings.Join(b.messages(), ", ")
}

func request(uri string, target interface{}) error {
	resp, err := http.Get(createURIString(uri))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		errors := map[string][]map[string]string{}
		if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
			return err
		}
		return BadRequest{errors: errors["errors"]}
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}

// Return a list of all known Licenses.
func All() (license.Licenses, error) {
	ret := license.Licenses{}
	return ret, request("licenses/", &ret)
}

// Return a list of all licenses which contain the keyword that was
// passed in.
func Tagged(keyword string) (license.Licenses, error) {
	ret := license.Licenses{}
	return ret, request("licenses/"+keyword, &ret)
}

func Get(id string) (license.License, error) {
	ret := license.License{}
	return ret, request("license/"+id, &ret)
}

// vim: foldmethod=marker
