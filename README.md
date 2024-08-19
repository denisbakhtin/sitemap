XML sitemap
=======

#Usage
```
  import (
    "github.com/denisbakhtin/sitemap"
    "time"
    ...
  )

  func CreateSitemap() {
    folder = "public_sitemap_folder"
    domain := "http://mydomain.com"
    now := time.Now()
    pages := make([]sitemap.Page, 1)

    //Home page
    pages = append(pages, sitemap.Page{
      Loc:        fmt.Sprintf("%s", domain),
      LastMod:    now,
      Changefreq: sitemap.Daily,
      Priority:   1,
    })

    //more pages
    published := models.GetPublishedPages() //get slice of pages
    for i := range published {
      pages = append(pages, sitemap.Page{
        Loc:        fmt.Sprintf("%s/pages/%d", domain, published[i].Id), //page url
        LastMod:    published[i].UpdatedAt, //page modification timestamp (time.Time)
        Changefreq: sitemap.Weekly, //or "hourly", "daily", ...
        Priority:   0.8,
      })
    }
    if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), pages); err != nil {
      log.Error(err)
      return
    }
    if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
      log.Error(err)
      return
    }
    //done
  }

```
For periodical sitemap generation pair it with something like this: https://github.com/jasonlvhit/gocron
