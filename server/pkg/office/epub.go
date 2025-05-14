package office

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"os"
	"path"
)

const containerPath = "META-INF/container.xml"

func epub2txt(filename string) (string, error) {
	text := ""
	rc, _ := OpenReader(filename)
	book := rc.Rootfiles[0]

	// Print book title.
	title := "Title: " + book.Title + "\n"
	// List the IDs of files in the book's spine.
	for _, item := range book.Spine.Itemrefs {
		// item.ID was observed to be in one book: cover,titlepage,brief-toc,xpreface_001,xintroduction_001,xepigraph_001,xchapter_001
		reader2, err := item.Open()
		if err != nil {
			log.Printf("Error opening item %s: %v", item.ID, err)
			continue
		}
		itemText, err := ConvertHTML(reader2)
		if err != nil {
			log.Printf("Convert HTML item %s: %v", item.ID, err)
			continue
		}
		text += itemText
	}

	if text == "" {
		return "", nil
	}
	return title + text, nil
}

var (
	// ErrNoRootfile occurs when there are no rootfile entries found in
	// container.xml.
	ErrNoRootfile = errors.New("epub: no rootfile found in container")

	// ErrBadRootfile occurs when container.xml references a rootfile that does
	// not exist in the zip.
	ErrBadRootfile = errors.New("epub: container references non-existent rootfile")

	// ErrNoItemref occurrs when a content.opf contains a spine without any
	// itemref entries.
	ErrNoItemref = errors.New("epub: no itemrefs found in spine")

	// ErrBadItemref occurs when an itemref entry in content.opf references an
	// item that does not exist in the manifest.
	ErrBadItemref = errors.New("epub: itemref references non-existent item")

	// ErrBadManifest occurs when a manifest in content.opf references an item
	// that does not exist in the zip.
	ErrBadManifest = errors.New("epub: manifest references non-existent item")
)

// Reader represents a readable epub file.
type Reader struct {
	Container
	files map[string]*zip.File
}

// ReadCloser represents a readable epub file that can be closed.
type ReadCloser struct {
	Reader
	f *os.File
}

// Rootfile contains the location of a content.opf package file.
type Rootfile struct {
	FullPath string `xml:"full-path,attr"`
	Package
}

// Container serves as a directory of Rootfiles.
type Container struct {
	Rootfiles []*Rootfile `xml:"rootfiles>rootfile"`
}

// Package represents an epub content.opf file.
type Package struct {
	Metadata
	Manifest
	Spine
}

// Metadata contains publishing information about the epub.
type Metadata struct {
	Title       string `xml:"metadata>title"`
	Language    string `xml:"metadata>language"`
	Identifier  string `xml:"metadata>idenifier"`
	Creator     string `xml:"metadata>creator"`
	Contributor string `xml:"metadata>contributor"`
	Publisher   string `xml:"metadata>publisher"`
	Subject     string `xml:"metadata>subject"`
	Description string `xml:"metadata>description"`
	Event       []struct {
		Name string `xml:"event,attr"`
		Date string `xml:",innerxml"`
	} `xml:"metadata>date"`
	Type     string `xml:"metadata>type"`
	Format   string `xml:"metadata>format"`
	Source   string `xml:"metadata>source"`
	Relation string `xml:"metadata>relation"`
	Coverage string `xml:"metadata>coverage"`
	Rights   string `xml:"metadata>rights"`
}

// Manifest lists every file that is part of the epub.
type Manifest struct {
	Items []Item `xml:"manifest>item"`
}

// Item represents a file stored in the epub.
type Item struct {
	ID        string `xml:"id,attr"`
	HREF      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
	f         *zip.File
}

// Spine defines the reading order of the epub documents.
type Spine struct {
	Itemrefs []Itemref `xml:"spine>itemref"`
}

// Itemref points to an Item.
type Itemref struct {
	IDREF string `xml:"idref,attr"`
	*Item
}

// OpenReader will open the epub file specified by name and return a
// ReadCloser.
func OpenReader(name string) (*ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	rc := new(ReadCloser)
	rc.f = f

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	z, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, err
	}

	if err = rc.init(z); err != nil {
		return nil, err
	}

	return rc, nil
}

// NewReader returns a new Reader reading from ra, which is assumed to have the
// given size in bytes.
func NewReader(ra io.ReaderAt, size int64) (*Reader, error) {
	z, err := zip.NewReader(ra, size)
	if err != nil {
		return nil, err
	}

	r := new(Reader)
	if err = r.init(z); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Reader) init(z *zip.Reader) error {
	// Create a file lookup table
	r.files = make(map[string]*zip.File)
	for _, f := range z.File {
		r.files[f.Name] = f
	}

	err := r.setContainer()
	if err != nil {
		return err
	}
	err = r.setPackages()
	if err != nil {
		return err
	}
	err = r.setItems()
	if err != nil {
		return err
	}

	return nil
}

// setContainer unmarshals the epub's container.xml file.
func (r *Reader) setContainer() error {
	f, err := r.files[containerPath].Open()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	_, err = io.Copy(&b, f)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(b.Bytes(), &r.Container)
	if err != nil {
		return err
	}

	if len(r.Container.Rootfiles) < 1 {
		return ErrNoRootfile
	}

	return nil
}

// setPackages unmarshal's each of the epub's content.opf files.
func (r *Reader) setPackages() error {
	for _, rf := range r.Container.Rootfiles {
		if r.files[rf.FullPath] == nil {
			return ErrBadRootfile
		}

		f, err := r.files[rf.FullPath].Open()
		if err != nil {
			return err
		}

		var b bytes.Buffer
		_, err = io.Copy(&b, f)
		if err != nil {
			return err
		}

		err = xml.Unmarshal(b.Bytes(), &rf.Package)
		if err != nil {
			return err
		}
	}

	return nil
}

// setItems associates Itemrefs with their respective Item and Items with
// their zip.File.
func (r *Reader) setItems() error {
	itemrefCount := 0
	for _, rf := range r.Container.Rootfiles {
		itemMap := make(map[string]*Item)
		for i := range rf.Manifest.Items {
			item := &rf.Manifest.Items[i]
			itemMap[item.ID] = item

			abs := path.Join(path.Dir(rf.FullPath), item.HREF)
			item.f = r.files[abs]
		}

		for i := range rf.Spine.Itemrefs {
			itemref := &rf.Spine.Itemrefs[i]
			itemref.Item = itemMap[itemref.IDREF]
			if itemref.Item == nil {
				return ErrBadItemref
			}
		}
		itemrefCount += len(rf.Spine.Itemrefs)
	}

	if itemrefCount < 1 {
		return ErrNoItemref
	}

	return nil
}

// Open returns a ReadCloser that provides access to the Items's contents.
// Multiple items may be read concurrently.
func (item *Item) Open() (r io.ReadCloser, err error) {
	if item.f == nil {
		return nil, ErrBadManifest
	}

	return item.f.Open()
}

// Close closes the epub file, rendering it unusable for I/O.
func (rc *ReadCloser) Close() {
	rc.f.Close()
}
