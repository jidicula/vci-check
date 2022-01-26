//     vci-check checks if a vaccine credential issuer is in the VCI Directory
// Copyright (C) 2022  Johanan Idicula
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jidicula/vci-check/atom"
	"github.com/jidicula/vci-check/checker"
)

// TODO: put web server logic in here
func main() {

	latest, _ := getLatestCommit()
	checkLatest(latest)
}

// checkLatest has an infinite loop to check the stored ID against latest ID.
func checkLatest(previousLatest string) {
	for {
		latest, err := getLatestCommit()
		if err != nil {
			continue
		}
		if latest != previousLatest {
			previousLatest = latest
			il, err := checker.NewIssuerList()
			if err != nil {
				fmt.Printf("%s", err)
			}
			// TODO: put into channel
			fmt.Println(il.ParticipatingIssuers[0].Name)
		}
		time.Sleep(5 * time.Second)
	}
}

// getLatestCommit parses the VCI commit Atom feed and returns the latest
// commit item.
func getLatestCommit() (string, error) {
	a := atom.Feed{}

	resp, err := http.Get("https://github.com/the-commons-project/vci-directory/commits/main.atom")
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status not OK, %d and %s", resp.StatusCode, body)
	}
	err = xml.Unmarshal(body, &a)
	return a.Entry[0].ID, err
}
