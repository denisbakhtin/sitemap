// Package sitemap provides structures and functions for web-site xml-sitemap management.
package sitemap

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	footer = `
	</urlset>`
	template = `
	<url>
	  <loc>%s</loc>
	  <lastmod>%s</lastmod>
	  <changefreq>%s</changefreq>
	  <priority>%.1f</priority>
	</url> 	`

	indexHeader = `<?xml version="1.0" encoding="UTF-8"?>
  <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexFooter = `
	</sitemapindex>`
	indexTemplate = `
	<sitemap>
		<loc>%s%s</loc>
		<lastmod>%s</lastmod>
	</sitemap>`
)

// Frequency represents a web-page refresh rate
type Frequency int

// Predefined constants for page refresh frequency.
// The value "always" should be used to describe documents that change each time they are accessed. The value "never" should be used to describe archived URLs.
// Please note that the value of this tag is considered a hint and not a command.
const (
	Always Frequency = iota
	Hourly
	Daily
	Weekly
	Monthly
	Yearly
	Never
)

func (f Frequency) String() string {
	return [...]string{"always", "hourly", "daily", "weekly", "monthly", "yearly", "never"}[f]
}

// Page represents a web-page record in a sitemap file.
// Loc is an absolute page URL.
// LastMod is a timestamp of page's last modification.
// Changefreq is how often page'c content is updated on server. Use constants Always, Hourly,... Never to set this.
// Priority is an arbitrary decimal number from 0 to 1 (ex: 0.7) that sets an importance priority for a search engine crawler.
type Page struct {
	Loc        string
	LastMod    time.Time
	Changefreq Frequency
	Priority   float32
}

// String returns a string representation of Page using a prefefined template
func (page *Page) String() string {
	return fmt.Sprintf(template, page.Loc, page.LastMod.Format("2006-01-02"), page.Changefreq, page.Priority)
}

// Sitemap creates an .xml.gz sitemap file with name f for a list of pages.
func SiteMap(f string, pages []Page) error {
	if !strings.HasSuffix(f, ".xml.gz") {
		return fmt.Errorf("not an .xml.gz filename provided")
	}

	var buffer bytes.Buffer
	buffer.WriteString(header)
	for _, page := range pages {
		_, err := buffer.WriteString(page.String())
		if err != nil {
			return err
		}
	}
	fo, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fo.Close()
	buffer.WriteString(footer)

	zip := gzip.NewWriter(fo)
	defer zip.Close()
	_, err = zip.Write(buffer.Bytes())
	return err
}

// SiteMapIndex creates an .xml index file named indexFile for all .xml.gz sitemap files in folder. Baseurl is a base web url for this folder.
func SiteMapIndex(folder, indexFile, baseurl string) error {
	var buffer bytes.Buffer
	buffer.WriteString(indexHeader)
	fs, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") {
			info, err := f.Info()
			if err != nil {
				return err
			}
			s := fmt.Sprintf(indexTemplate, baseurl, f.Name(), info.ModTime().Format("2006-01-02"))
			buffer.WriteString(s)
		}
	}
	buffer.WriteString(indexFooter)

	fo, err := os.Create(path.Join(folder, indexFile))
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = fo.Write(buffer.Bytes())
	return err
}
