package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Modfile struct {
	ProjectFileID uint64
	FileName      string
}

func ParseModfile(jsonVal string) (file Modfile) {
	json.Unmarshal([]byte(jsonVal), &file)
	return
}

type Mod struct {
	Game  string
	Addon string

	Created uint64
	Updated uint64
	Install uint64

	Installed Modfile
	Newest    Modfile
}

func (m Mod) CurseURL() string {
	return fmt.Sprintf("http://www.curseforge.com/%s/%s", m.Game, m.Addon)
}

func (m Mod) DownloadURL(file Modfile) string {
	return fmt.Sprintf("http://www.curseforge.com/%s/%s/download/%d/file", m.Game, m.Addon, file.ProjectFileID)
}

func (m Mod) cacheDirectory() string {
	dir := path.Join(addoncache, m.Game, m.Addon)
	os.MkdirAll(dir, os.ModePerm)
	return dir
}

func (m *Mod) Update() {
	// pull the info page
	doc, err := goquery.NewDocument(m.CurseURL())
	if err != nil {
		return
	}

	// creation/update times
	times := doc.Find("abbr")
	created, _ := times.Eq(1).Attr("data-epoch")
	m.Created, _ = strconv.ParseUint(created, 10, 64)
	updated, _ := times.Eq(0).Attr("data-epoch")
	m.Updated, _ = strconv.ParseUint(updated, 10, 64)

	// pull the release list
	files, err := goquery.NewDocument(m.CurseURL() + "/files")
	if err != nil {
		return
	}

	modfileJson, _ := files.Find("a.button--download").Eq(1).Attr("data-action-value")
	m.Newest = ParseModfile(modfileJson)
}

func (m *Mod) InstallTo() {
	// download zip from mod, install to directory by just xtract
	resp, err := http.Get(m.DownloadURL(m.Newest))
	if err != nil {
		return
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ioutil.WriteFile(path.Join(m.cacheDirectory(), m.Newest.FileName+".zip"), dat, os.ModePerm)
	rto := bytes.NewReader(dat)

	zr, err := zip.NewReader(rto, resp.ContentLength)
	if err != nil {
		return
	}

	for _, file := range zr.File {
		zfr, err := file.Open()
		if err != nil {
			return
		}

		dest := adir + file.Name
		os.MkdirAll(path.Dir(dest), os.ModePerm)
		d, _ := os.Create(dest)
		io.Copy(d, zfr)
		d.Close()
	}
	m.Install = m.Updated
	m.Installed = m.Newest
	log.Printf("finished install (%s to version %s)", m.Addon, m.Installed.FileName)
}

func (m *Mod) Uninstall() {
	// do this, added the things
	// also check docs for return (probably error)
	log.Print("possibly finished uninstall")
}

func (m *Mod) UpdateI() {
	if m.Install < m.Updated {
		m.InstallTo()
	}
}

func GetMod(game, addon string) *Mod {
	r := &Mod{Game: game, Addon: addon}
	r.Update()
	return r
}
