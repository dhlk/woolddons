package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// based on github.com/paralin/CURSETools
type Mod struct {
	Created uint64
	Updated uint64
	Install uint64
	Version string

	CurseURL    string
	ProjectURL  string
	DownloadURL string
}

func (m *Mod) Update() {
	doc, err := goquery.NewDocument(m.CurseURL)
	if err != nil {
		return
	}

	body := doc.Find("div.main-info")
	details := body.Find("ul.details-list")

	// creation/update times
	details.Find("li.updated").Each(
		func(i int, s *goquery.Selection) {
			update, _ := s.Children().First().Attr("data-epoch")
			if strings.HasPrefix(s.Text(), "Updated") {
				m.Updated, _ = strconv.ParseUint(update, 10, 64)
			} else {
				m.Created, _ = strconv.ParseUint(update, 10, 64)
			}
		})

	// version
	nft := details.Find("li.newest-file").Text()
	vid := strings.Index(nft, ":") + 2
	if vid < len(nft) {
		m.Version = nft[vid:]
	} else {
		m.Version = "UNK - TODO mod.go:57" // TODO
	}

	// curseforge url
	m.ProjectURL, _ = details.Find(".curseforge").Find("a").Attr("href")

	// download url
	doc, err = goquery.NewDocument(m.CurseURL + "/download")
	if err != nil {
		m.DownloadURL = err.Error()
	}
	m.DownloadURL, _ = doc.Find("div.countdown").Find("a").Attr("data-href")
}

func (m *Mod) InstallTo() {
	// download zip from mod, install to directory by just xtract
	resp, err := http.Get(m.DownloadURL)
	if err != nil {
		return
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ioutil.WriteFile(addoncache+"/"+path.Base(resp.Request.URL.Path), dat, os.ModePerm)
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
	fmt.Printf("finished install (%s to version %s)\n", m.CurseURL, m.Version)
}

func (m *Mod) Uninstall() {
	// do this, added the things
	// also check docs for return (probably error)
	fmt.Printf("possibly finished uninstall\n")
}

func (m *Mod) UpdateI() {
	if m.Install < m.Updated {
		m.InstallTo()
	}
}

func ParseURL(url string) *Mod {
	r := &Mod{CurseURL: url}
	r.Update()
	return r
}

func GetMod(game, addon string) *Mod {
	return ParseURL(fmt.Sprintf("http://www.curse.com/%s/%s", game, addon))
}
