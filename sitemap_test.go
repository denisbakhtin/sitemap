package sitemap

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

func TestPageString(t *testing.T) {
	lastMod := time.Date(2024, 12, 1, 1, 1, 1, 1, time.Local)
	page := Page{Loc: "/111.html", LastMod: lastMod, Changefreq: Daily, Priority: 0.2}
	got := page.String()
	want := `
	<url>
	  <loc>/111.html</loc>
	  <lastmod>2024-12-01</lastmod>
	  <changefreq>daily</changefreq>
	  <priority>0.2</priority>
	</url> 	`

	if got != want {
		t.Errorf("Page.String()  = %q; want %q", got, want)
	}
}

func TestSiteMap(t *testing.T) {
	pages := []Page{{Loc: "http://example.com/1.html", Changefreq: Hourly, Priority: 0.7, LastMod: time.Now()}}
	if err := SiteMap("some.xml", pages); err == nil {
		t.Errorf("SiteMap should accept only .xml.gz file names")
	}

	file, err := os.CreateTemp("", "*.xml.gz")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer file.Close()
	if err := SiteMap(file.Name(), pages); err != nil {
		t.Errorf(err.Error())
	}

	zip, err := gzip.NewReader(file)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer zip.Close()
	bb, err := io.ReadAll(zip)

	if err != nil {
		t.Errorf(err.Error())
	}
	s := string(bb)
	if s != header+pages[0].String()+footer {
		t.Errorf("Wrong sitemap file content")
	}
}

func TestSiteMapIndex(t *testing.T) {
	dir, err := os.MkdirTemp("", "*sitemaps")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer os.RemoveAll(dir)
	file, err := os.CreateTemp(dir, "*.xml.gz")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer file.Close()
	if err := SiteMapIndex(dir, "sitemap_index.xml", "/sitemaps/"); err != nil {
		t.Errorf(err.Error())
	}
	index, err := os.ReadFile(path.Join(dir, "sitemap_index.xml"))
	if err != nil {
		t.Errorf(err.Error())
	}
	s := string(index)
	info, err := os.Stat(file.Name())
	if err != nil {
		t.Errorf(err.Error())
	}
	want := indexHeader + fmt.Sprintf(indexTemplate, "/sitemaps/", filepath.Base(file.Name()), info.ModTime().Format("2006-01-02")) + indexFooter
	if s != want {
		t.Errorf("Wrong sitemap index contents:\nGot: %q\nWant: %q\n", s, want)
	}
}
