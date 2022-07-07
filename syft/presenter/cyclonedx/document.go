package cyclonedx

import (
	"encoding/xml"

	"github.com/zj1244/syft/internal"
	"github.com/zj1244/syft/internal/version"
	"github.com/zj1244/syft/syft/pkg"
	"github.com/zj1244/syft/syft/source"
	"github.com/google/uuid"
)

// Source: https://github.com/CycloneDX/specification

// Document represents a CycloneDX BOM Document.
type Document struct {
	XMLName       xml.Name       `xml:"bom"`
	XMLNs         string         `xml:"xmlns,attr"`
	Version       int            `xml:"version,attr"`
	SerialNumber  string         `xml:"serialNumber,attr"`
	BomDescriptor *BomDescriptor `xml:"metadata"`             // The BOM descriptor extension
	Components    []Component    `xml:"components>component"` // The BOM contents
}

// NewDocumentFromCatalog returns a CycloneDX Document object populated with the catalog contents.
func NewDocument(catalog *pkg.Catalog, srcMetadata source.Metadata) Document {
	versionInfo := version.FromBuild()

	doc := Document{
		XMLNs:         "http://cyclonedx.org/schema/bom/1.2",
		Version:       1,
		SerialNumber:  uuid.New().URN(),
		BomDescriptor: NewBomDescriptor(internal.ApplicationName, versionInfo.Version, srcMetadata),
	}

	// attach components
	for p := range catalog.Enumerate() {
		component := Component{
			Type:       "library", // TODO: this is not accurate
			Name:       p.Name,
			Version:    p.Version,
			PackageURL: p.PURL,
		}
		var licenses []License
		for _, licenseName := range p.Licenses {
			licenses = append(licenses, License{
				Name: licenseName,
			})
		}
		if len(licenses) > 0 {
			component.Licenses = &licenses
		}
		doc.Components = append(doc.Components, component)
	}

	return doc
}
