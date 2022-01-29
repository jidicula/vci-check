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
	"log"
	"net/http"

	"github.com/jidicula/vci-check/atom"
	"github.com/jidicula/vci-check/checker"
)

func main() {

	ilCh := make(chan checker.IssuerList)
	handleCheck := makeHandleCheck(ilCh)

	latest, _ := getLatestCommit()
	go checkLatest(latest, ilCh)

	http.HandleFunc("/", handleCheck)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func makeHandleCheck(ilCh <-chan checker.IssuerList) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Accept") != "application/json" {
			w.Header().Set("Accept", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// serveraddress:port?iss=<url>
		issURL := r.URL.Query().Get("iss")
		if issURL == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		il := <-ilCh
		response := fmt.Sprintf(`{"message": %t}`, il.IsTrusted(issURL))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		x, err := w.Write([]byte(response))
		if err != nil {
			fmt.Printf("%s, %d bytes written", err, x)
		}
		fmt.Printf("done\n")
	}
}

// checkLatest has an infinite loop to check the stored ID against latest ID.
func checkLatest(previousLatest string, ilCh chan<- checker.IssuerList) {
	// get issuerList on function init
	il, err := checker.NewIssuerList()
	if err != nil {
		fmt.Printf("%s", err)
	}
	for {
		latest, err := getLatestCommit()
		if err != nil {
			continue
		}
		if latest != previousLatest {
			previousLatest = latest
			// update issuerList if there's new commits
			il, err = checker.NewIssuerList()
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Println(il.ParticipatingIssuers[0].Name)
		}
		// Always write issuerList into channel on each iteration.
		// Blocks until ilCh is read from in handleCheck
		ilCh <- il
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
