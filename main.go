package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const addoninfo, addoncache string = "addoncache/addoncache.json", "addoncache"

func cleanname(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "-", -1))
}

type AddonList map[string]*Mod

func (al AddonList) Save() {
	dat, err := json.Marshal(al)
	if err != nil {
		//
		return
	}
	ioutil.WriteFile(addoninfo, dat, os.ModePerm)
}

func (al AddonList) Load() {
	dat, err := ioutil.ReadFile(addoninfo)
	if err != nil {
		//
		return
	}
	json.Unmarshal(dat, &al)
}

func (al AddonList) Add(game, name string) {
	mod := GetMod(game, name)
	if mod != nil {
		al[name] = mod
	}
}

var port, adir, game string

func main() {
	flag.StringVar(&port, "addr", ":8080", "webserver port")
	flag.StringVar(&adir, "adir", "addons", "directory to install addons into")
	flag.StringVar(&game, "game", "wow/addons", "game prefix (wildstar/ws-addons for wildstar...)")
	flag.Parse()
	adir = adir + "/"

	addons := AddonList{}
	addons.Load()

	http.HandleFunc(
		"/refresh",
		func(w http.ResponseWriter, req *http.Request) {
			for _, mod := range addons {
				mod.Update()
				mod.UpdateI()
				addons.Save()
			}
			http.Redirect(w, req, "/", http.StatusFound)
		})

	http.HandleFunc(
		"/act",
		func(w http.ResponseWriter, req *http.Request) {
			// add,remove
			action, addon := req.PostFormValue("action"), cleanname(req.PostFormValue("addon"))
			if addon == "" {
				return
			}
			if action == "add" {
				addons.Add(game, addon)
				addons[addon].InstallTo()
			} else if action == "remove" {
				addons[addon].Uninstall()
				delete(addons, addon)
			}
			addons.Save()
			http.Redirect(w, req, "/", http.StatusFound)
		})

	// home
	tmpl := template.Must(template.Must(template.New("base").Parse(page)).Parse(`{{template "page" .}}`))
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, req *http.Request) {
			tmpl.Execute(w, addons)
		})

	http.HandleFunc(
		"/style.css",
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, style)
		})

	http.ListenAndServe(port, nil)
}
