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

type Frequency int

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

type Item struct {
	Loc        string
	LastMod    time.Time
	Changefreq Frequency
	Priority   float32
}

func (item *Item) String() string {
	return fmt.Sprintf(template, item.Loc, item.LastMod.Format("2006-01-02"), item.Changefreq, item.Priority)
}

func SiteMap(f string, items []Item) error {
	var buffer bytes.Buffer
	buffer.WriteString(header)
	for _, item := range items {
		_, err := buffer.WriteString(item.String())
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

func SiteMapIndex(folder, indexFile, baseurl string) error {
	var buffer bytes.Buffer
	buffer.WriteString(indexHeader)
	fs, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") {
			s := fmt.Sprintf(indexTemplate, baseurl, f.Name(), time.Now().Format("2006-01-02"))
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
