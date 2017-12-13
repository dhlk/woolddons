package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"

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

	mu *sync.Mutex
}

func (m *Mod) lock() {
	if m.mu == nil {
		m.mu = &sync.Mutex{}
	}
	m.mu.Lock()
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
	m.lock()
	defer m.mu.Unlock()

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
	m.lock()
	defer m.mu.Unlock()

	zipName := path.Join(m.cacheDirectory(), m.Newest.FileName+".zip")

	err := ZipDownload(m.DownloadURL(m.Newest), zipName)
	if err != nil {
		return
	}

	err = ZipInstall(zipName, adir)
	if err != nil {
		return
	}

	m.Install = m.Updated
	m.Installed = m.Newest

	log.Printf("finished install (%s to version %s)", m.Addon, m.Installed.FileName)
}

func (m *Mod) Uninstall() {
	m.lock()
	defer m.mu.Unlock()

	zipName := path.Join(m.cacheDirectory(), m.Installed.FileName+".zip")

	// if you foolishly deleted the cached copy...
	if _, err := os.Stat(zipName); os.IsNotExist(err) {
		err = ZipDownload(m.DownloadURL(m.Installed), zipName)
		if err != nil {
			return
		}
	}

	err := ZipUninstall(zipName, adir)
	if err != nil {
		return
	}

	defer func() {
		m.Install = 0
		m.Installed = Modfile{}
	}()

	log.Printf("finished uninstall (%s from version %s)", m.Addon, m.Installed.FileName)
}

func (m *Mod) UpdateI() {
	if m.Install < m.Updated {
		m.Uninstall()
		m.InstallTo()
	}
}

func GetMod(game, addon string) *Mod {
	r := &Mod{Game: game, Addon: addon}
	r.Update()
	return r
}
