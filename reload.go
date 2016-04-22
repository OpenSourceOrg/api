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
	"os"
	"os/signal"
	"syscall"

	"github.com/opensourceorg/api/license"
)

type Blobs struct {
	Licenses      license.Licenses
	LicenseIdMap  map[string]license.License
	LicenseTagMap map[string][]license.License
}

func loadBlob(path string, blob *Blobs) error {
	licenses, err := license.LoadLicensesFiles(path)
	if err != nil {
		return err
	}
	licenseIdMap := licenses.GetIdMap()
	licenseTagMap := licenses.GetTagMap()

	blob.Licenses = licenses
	blob.LicenseIdMap = licenseIdMap
	blob.LicenseTagMap = licenseTagMap

	return nil
}

func Reloader(file string, target *Blobs) {
	if err := loadBlob(file, target); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	for _ = range c {
		if err := loadBlob(file, target); err != nil {
			panic(err)
		}
	}
}

// vim: foldmethod=marker
