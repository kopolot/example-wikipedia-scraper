package scraper

import "encoding/xml"

type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type SitemapIndex struct {
	XMLName  xml.Name   `xml:"sitemapindex"`
	Sitemaps []*Sitemap `xml:"sitemap"`
}

type SitemapURLSet struct {
	URLs []*SitemapUrl `xml:"url"`
}

type SitemapUrl struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}
