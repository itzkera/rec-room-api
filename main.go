
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)


type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type PlayerProgression struct {
	ID    int64 `json:"id"`
	Level int   `json:"level"`
	XP    int   `json:"xp"`
}

type PlayerEvent struct {
	ID        int64  `json:"id"`
	CreatorID int64  `json:"creatorId"`
	RoomID    int64  `json:"roomId"`
	Name      string `json:"name"`
}

type Image struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Owner  int64  `json:"owner"`
	Cheers int    `json:"cheers"`
}

var accounts = []Account{
	{1, "kera", "kera"},
}

var progressions = []PlayerProgression{
	{1, 10, 1500},
	{2, 25, 9000},
}

var events = []PlayerEvent{
	{1, 1, 1, "Launch Party"},
	{2, 2, 2, "PvP Tourney"},
}

var images = []Image{
	{1, "img1", 1, 5},
}


func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":8080", nil)
}



func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/progression/bulk" {
		handleProgressionBulk(w, r)
		return
	}

	if path == "/events" {
		writeJSON(w, events)
		return
	}

	if strings.HasPrefix(path, "/events/creator/") {
		id := parseID(path, "/events/creator/")
		filterEvents(w, func(e PlayerEvent) bool {
			return e.CreatorID == id
		})
		return
	}

	if strings.HasPrefix(path, "/events/room/") {
		id := parseID(path, "/events/room/")
		filterEvents(w, func(e PlayerEvent) bool {
			return e.RoomID == id
		})
		return
	}

	if strings.HasPrefix(path, "/events/search") {
		q := r.URL.Query().Get("query")
		filterEvents(w, func(e PlayerEvent) bool {
			return strings.Contains(strings.ToLower(e.Name), strings.ToLower(q))
		})
		return
	}

	if strings.HasPrefix(path, "/events/") {
		id := parseID(path, "/events/")
		for _, e := range events {
			if e.ID == id {
				writeJSON(w, e)
				return
			}
		}
	}

	if path == "/isinfluencer" {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		isInfluencer := id == 1 
		writeJSON(w, map[string]bool{
			"isInfluencer": isInfluencer,
		})
		return
	}
	if strings.HasPrefix(path, "/images/") {
		handleImages(w, r)
		return
	}

	if path == "/accounts" {
		username := r.URL.Query().Get("username")
		for _, a := range accounts {
			if a.Username == username {
				writeJSON(w, a)
				return
			}
		}
		writeJSON(w, nil)
		return
	}

	if strings.HasPrefix(path, "/accounts/search") {
		q := r.URL.Query().Get("name")
		filterAccounts(w, func(a Account) bool {
			return strings.Contains(strings.ToLower(a.Username), strings.ToLower(q))
		})
		return
	}

	if path == "/accounts/bulk" {
		handleAccountsBulk(w, r)
		return
	}

	if strings.HasPrefix(path, "/accounts/") {
		handleAccountByID(w, r)
		return
	}

	http.NotFound(w, r)
}


func handleProgressionBulk(w http.ResponseWriter, r *http.Request) {
	var result []PlayerProgression

	if r.Method == "GET" {
		ids := r.URL.Query()["id"]
		for _, idStr := range ids {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			for _, p := range progressions {
				if p.ID == id {
					result = append(result, p)
				}
			}
		}
	} else {
		r.ParseForm()
		for _, idStr := range r.Form["id"] {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			for _, p := range progressions {
				if p.ID == id {
					result = append(result, p)
				}
			}
		}
	}

	writeJSON(w, result)
}

func handleImages(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/public/images")
	// like apim.rec.net/public/images <3

	if strings.Contains(path, "/cheers") {
		id := parseID(path, "")
		for _, i := range images {
			if i.ID == id {
				writeJSON(w, i.Cheers)
				return
			}
		}
	}

	if strings.Contains(path, "/comments") {
		writeJSON(w, []string{"Nice!", "Cool!"})
		return
	}

	id, _ := strconv.ParseInt(strings.Split(path, "/")[0], 10, 64)
	for _, i := range images {
		if i.ID == id {
			writeJSON(w, i)
			return
		}
	}
}

func handleAccountsBulk(w http.ResponseWriter, r *http.Request) {
	var result []Account

	if r.Method == "GET" {
		for _, idStr := range r.URL.Query()["id"] {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			for _, a := range accounts {
				if a.ID == id {
					result = append(result, a)
				}
			}
		}
	} else {
		r.ParseForm()
		for _, idStr := range r.Form["id"] {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			for _, a := range accounts {
				if a.ID == id {
					result = append(result, a)
				}
			}
		}
	}

	writeJSON(w, result)
}

func handleAccountByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")

	if strings.HasSuffix(path, "/bio") {
		id := parseID(path, "")
		for _, a := range accounts {
			if a.ID == id {
				writeJSON(w, a.Bio)
				return
			}
		}
		return
	}

	id, _ := strconv.ParseInt(path, 10, 64)
	for _, a := range accounts {
		if a.ID == id {
			writeJSON(w, a)
			return
		}
	}
}


func parseID(path, prefix string) int64 {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.Split(idStr, "/")[0]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

func filterEvents(w http.ResponseWriter, fn func(PlayerEvent) bool) {
	var result []PlayerEvent
	for _, e := range events {
		if fn(e) {
			result = append(result, e)
		}
	}
	writeJSON(w, result)
}

func filterAccounts(w http.ResponseWriter, fn func(Account) bool) {
	var result []Account
	for _, a := range accounts {
		if fn(a) {
			result = append(result, a)
		}
	}
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}