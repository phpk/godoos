package xlsx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"sync"
)

// XlsxFile defines a populated XLSX file struct.
type XlsxFile struct {
	Sheets []string

	sheetFiles    map[string]*zip.File
	sharedStrings []string
	dateStyles    map[int]bool

	doneCh chan struct{} // doneCh serves as a signal to abort unfinished operations.
}

// XlsxFileCloser wraps XlsxFile to be able to close an open file
type XlsxFileCloser struct {
	zipReadCloser *zip.ReadCloser
	XlsxFile

	once sync.Once // once performs actions exactly once, e.g. closing a channel.
}

// getFileForName finds and returns a *zip.File by it's display name from within an archive.
// If the file cannot be found, an error is returned.
func getFileForName(files []*zip.File, name string) (*zip.File, error) {
	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", name)
}

// readFile opens and reads the entire contents of a *zip.File into memory.
// If the file cannot be opened, or the data cannot be read, an error is returned.
func readFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return []byte{}, fmt.Errorf("unable to open file: %w", err)
	}
	defer rc.Close()

	buff := bytes.NewBuffer(nil)
	_, err = io.Copy(buff, rc)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to copy bytes: %w", err)
	}
	return buff.Bytes(), nil
}

// GetSheetFileForSheetName returns the sheet file associated with the sheet name.
// This is useful when you want to further process something out of the sheet, that this
// library does not handle. For example this is useful when trying to read the hyperlinks
// section of a sheet file; getting the sheet file enables you to read the XML directly.
func (xl *XlsxFileCloser) GetSheetFileForSheetName(sheetName string) *zip.File {
	sheetFile, _ := xl.sheetFiles[sheetName]
	return sheetFile
}

// Close closes the XlsxFile, rendering it unusable for I/O.
func (xl *XlsxFileCloser) Close() error {
	if xl == nil {
		return nil
	}
	xl.once.Do(func() { close(xl.doneCh) })
	return xl.zipReadCloser.Close()
}

// OpenFile takes the name of an XLSX file and returns a populated XlsxFile struct for it.
// If the file cannot be found, or key parts of the files contents are missing, an error
// is returned.
// Note that the file must be Close()-d when you are finished with it.
func OpenFile(filename string) (*XlsxFileCloser, error) {
	zipFile, err := zip.OpenReader(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file reader: %w", err)
	}

	x := XlsxFile{}
	if err := x.init(&zipFile.Reader); err != nil {
		zipFile.Close()
		return nil, fmt.Errorf("unable to initialise file: %w", err)
	}

	return &XlsxFileCloser{
		XlsxFile:      x,
		zipReadCloser: zipFile,
	}, nil
}

// OpenReaderZip takes the zip ReadCloser of an XLSX file and returns a populated XlsxFileCloser struct for it.
// If the file cannot be found, or key parts of the files contents are missing, an error
// is returned.
// Note that the file must be Close()-d when you are finished with it.
func OpenReaderZip(rc *zip.ReadCloser) (*XlsxFileCloser, error) {
	x := XlsxFile{}

	if err := x.init(&rc.Reader); err != nil {
		rc.Close()
		return nil, err
	}

	return &XlsxFileCloser{
		XlsxFile:      x,
		zipReadCloser: rc,
	}, nil
}

// NewReader takes bytes of Xlsx file and returns a populated XlsxFile struct for it.
// If the file cannot be found, or key parts of the files contents are missing, an error
// is returned.
func NewReader(xlsxBytes []byte) (*XlsxFile, error) {
	r, err := zip.NewReader(bytes.NewReader(xlsxBytes), int64(len(xlsxBytes)))
	if err != nil {
		return nil, fmt.Errorf("unable to create new reader: %w", err)
	}

	x := XlsxFile{}
	err = x.init(r)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise file: %w", err)
	}

	return &x, nil
}

// NewReaderZip takes zip reader of Xlsx file and returns a populated XlsxFile struct for it.
// If the file cannot be found, or key parts of the files contents are missing, an error
// is returned.
func NewReaderZip(r *zip.Reader) (*XlsxFile, error) {
	x := XlsxFile{}

	if err := x.init(r); err != nil {
		return nil, fmt.Errorf("unable to initialise file: %w", err)
	}

	return &x, nil
}

func (x *XlsxFile) init(zipReader *zip.Reader) error {
	sharedStrings, err := getSharedStrings(zipReader.File)
	if err != nil {
		return fmt.Errorf("unable to get shared strings: %w", err)
	}

	sheets, sheetFiles, err := getWorksheets(zipReader.File)
	if err != nil {
		return fmt.Errorf("unable to get worksheets: %w", err)
	}

	dateStyles, err := getDateFormatStyles(zipReader.File)
	if err != nil {
		return fmt.Errorf("unable to get date styles: %w", err)
	}

	x.sharedStrings = sharedStrings
	x.Sheets = sheets
	x.sheetFiles = *sheetFiles
	x.dateStyles = *dateStyles
	x.doneCh = make(chan struct{})

	return nil
}
