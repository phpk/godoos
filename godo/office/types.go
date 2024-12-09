package office

import (
	"regexp"
	"time"
)

const maxBytes = 1024 << 20 // 1GB
const ISO string = "2006-01-02T15:04:05"

var TAG_RE = regexp.MustCompile(`(<[^>]*>)+`)
var PARA_RE = regexp.MustCompile(`(</[a-z]:p>)+`)
var DEBUG bool = false

type Document struct {
	path           string
	RePath         string    `json:"path"`
	Filename       string    `json:"filename"`
	Title          string    `json:"title"`
	Subject        string    `json:"subject"`
	Creator        string    `json:"creator"`
	Keywords       string    `json:"keywords"`
	Description    string    `json:"description"`
	Lastmodifiedby string    `json:"lastModifiedBy"`
	Revision       string    `json:"revision"`
	Category       string    `json:"category"`
	Content        string    `json:"content"`
	Split          []string  `json:"split"`
	Modifytime     time.Time `json:"modified"`
	Createtime     time.Time `json:"created"`
	Accesstime     time.Time `json:"accessed"`
	Size           int       `json:"size"`
}

type DocReader func(string) (string, error)
type XMLContent struct {
	Title          string `xml:"title"`
	Subject        string `xml:"subject"`
	Creator        string `xml:"creator"`
	Keywords       string `xml:"keywords"`
	Description    string `xml:"description"`
	LastModifiedBy string `xml:"lastModifiedBy"`
	Revision       string `xml:"revision"`
	Created        string `xml:"created"`
	Modified       string `xml:"modified"`
	Category       string `xml:"category"`
}
