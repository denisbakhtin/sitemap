package sitemap

import (
	"testing"
	"time"
)

func TestItemString(t *testing.T) {
	lastMod := time.Date(2024, 12, 1, 1, 1, 1, 1, time.Local)
	item := Page{Loc: "/111.html", LastMod: lastMod, Changefreq: Daily, Priority: 0.2}
	got := item.String()
	want := `
	<url>
	  <loc>/111.html</loc>
	  <lastmod>2024-12-01</lastmod>
	  <changefreq>daily</changefreq>
	  <priority>0.2</priority>
	</url> 	`

	if got != want {
		t.Errorf("Item.String()  = %q; want %q", got, want)
	}
}
