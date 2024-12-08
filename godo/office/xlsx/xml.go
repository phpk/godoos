package xlsx

import (
	"encoding/xml"
	"fmt"
)

func getCharData(d *xml.Decoder) (string, error) {
	tok, err := d.Token()
	if err != nil {
		return "", fmt.Errorf("unable to get raw token: %w", err)
	}

	cdata, ok := tok.(xml.CharData)
	if !ok {
		// Valid for no chardata to be present
		return "", nil
	}

	return string(cdata), nil
}
