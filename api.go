/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@opensource.org>, 2015-2016
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>. }}} */

package main

import (
	"encoding/json"
	"net/http"
	"os"
)

//
func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

//
func writeError(w http.ResponseWriter, message string, code int) error {
	return writeJSON(w, map[string][]map[string]string{
		"errors": []map[string]string{
			map[string]string{"message": message},
		},
	}, code)
}

func main() {
	mux := http.NewServeMux()
	blob := Blobs{}

	go Reloader(os.Args[1], &blob)

	licensesEndpoint := "/licenses/"
	mux.HandleFunc(licensesEndpoint, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == licensesEndpoint {
			writeJSON(w, blob.Licenses, 200)
			return
		}
		path := req.URL.Path[len(licensesEndpoint):]
		if licenses, ok := blob.LicenseTagMap[path]; ok {
			writeJSON(w, licenses, 200)
			return
		}
		writeError(w, "Unknown tag", 404)
	})

	licenseEndpoint := "/license/"
	mux.HandleFunc(licenseEndpoint, func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path[len(licenseEndpoint):]
		if license, ok := blob.LicenseIdMap[path]; ok {
			writeJSON(w, license, 200)
			return
		}
		writeError(w, "Unknown license", 404)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		writeError(w, "No such page", 404)
	})

	http.ListenAndServe(":8000", mux)
}

// vim: foldmethod=marker
