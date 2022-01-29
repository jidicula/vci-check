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
	"fmt"
	"log"
	"net/http"

	"github.com/jidicula/vci-check/checker"
)

func main() {

	ilCh := make(chan checker.IssuerList)
	handleCheck := makeHandleCheck(ilCh)

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

		il, err := checker.NewIssuerList()
		if err != nil {
			http.Error(w, "Error retrieving VCI issuer list", http.StatusInternalServerError)
			return
		}
		response := fmt.Sprintf(`{"message": %t}`, il.IsTrusted(issURL))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		x, err := w.Write([]byte(response))
		if err != nil {
			fmt.Printf("%s, %d bytes written", err, x)
		}
	}
}
