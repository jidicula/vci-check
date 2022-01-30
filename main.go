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
	flag "github.com/spf13/pflag"
)

func main() {
	port := flag.String("port", "8080", "Port to listen to")
	flag.Parse()
	log.Printf("listening on %s", *port)

	http.HandleFunc("/", checkHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func checkHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		log.Printf("incorrect request method: %s", r.Method)
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "expect method GET at /?iss=<url>", http.StatusMethodNotAllowed)
		return
	}

	// serveraddress:port?iss=<url>
	issURL := r.URL.Query().Get("iss")
	if issURL == "" {
		msg := "No issuer URL provided"
		log.Print(msg)
		http.Error(w, fmt.Sprintf(`{"message": "%s"}`, msg), http.StatusBadRequest)
		return
	}

	il, err := checker.NewIssuerList()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, `{"message": "Error retrieving VCI issuer list"}`, http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf(`{"message": %t}`, il.IsTrusted(issURL))
	log.Printf("issURL %s : %s", issURL, response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	x, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("%s, %d bytes written", err, x)
	}
}
