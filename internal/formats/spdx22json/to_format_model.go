package spdx22json

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/zj1244/syft/syft/sbom"

	"github.com/zj1244/syft/internal"
	"github.com/zj1244/syft/internal/formats/common/spdxhelpers"
	"github.com/zj1244/syft/internal/formats/spdx22json/model"
	"github.com/zj1244/syft/internal/spdxlicense"
	"github.com/zj1244/syft/internal/version"
	"github.com/zj1244/syft/syft/pkg"
	"github.com/zj1244/syft/syft/source"
	"github.com/google/uuid"
)

// toFormatModel creates and populates a new JSON document struct that follows the SPDX 2.2 spec from the given cataloging results.
func toFormatModel(s sbom.SBOM) model.Document {
	name := documentName(s.Source)
	packages, files, relationships := extractFromCatalog(s.Artifacts.PackageCatalog)

	return model.Document{
		Element: model.Element{
			SPDXID: model.ElementID("DOCUMENT").String(),
			Name:   name,
		},
		SPDXVersion: model.Version,
		CreationInfo: model.CreationInfo{
			Created: time.Now().UTC(),
			Creators: []string{
				// note: key-value format derived from the JSON example document examples: https://github.com/spdx/spdx-spec/blob/v2.2/examples/SPDXJSONExample-v2.2.spdx.json
				"Organization: Anchore, Inc",
				"Tool: " + internal.ApplicationName + "-" + version.FromBuild().Version,
			},
			LicenseListVersion: spdxlicense.Version,
		},
		DataLicense:       "CC0-1.0",
		DocumentNamespace: documentNamespace(name, s.Source),
		Packages:          packages,
		Files:             files,
		Relationships:     relationships,
	}
}

func documentName(srcMetadata source.Metadata) string {
	switch srcMetadata.Scheme {
	case source.ImageScheme:
		return cleanSPDXName(srcMetadata.ImageMetadata.UserInput)
	case source.DirectoryScheme:
		return cleanSPDXName(srcMetadata.Path)
	}

	// TODO: is this alright?
	return uuid.Must(uuid.NewRandom()).String()
}

func documentNamespace(name string, srcMetadata source.Metadata) string {
	input := "unknown-source-type"
	switch srcMetadata.Scheme {
	case source.ImageScheme:
		input = "image"
	case source.DirectoryScheme:
		input = "dir"
	}

	uniqueID := uuid.Must(uuid.NewRandom())
	identifier := path.Join(input, uniqueID.String())
	if name != "." {
		identifier = path.Join(input, fmt.Sprintf("%s-%s", name, uniqueID.String()))
	}

	return path.Join(anchoreNamespace, identifier)
}

func extractFromCatalog(catalog *pkg.Catalog) ([]model.Package, []model.File, []model.Relationship) {
	packages := make([]model.Package, 0)
	relationships := make([]model.Relationship, 0)
	files := make([]model.File, 0)

	for _, p := range catalog.Sorted() {
		license := spdxhelpers.License(p)
		packageSpdxID := model.ElementID(fmt.Sprintf("Package-%+v-%s-%s", p.Type, p.Name, p.Version)).String()

		packageFiles, fileIDs, packageFileRelationships := spdxhelpers.Files(packageSpdxID, p)
		files = append(files, packageFiles...)

		relationships = append(relationships, packageFileRelationships...)

		// note: the license concluded and declared should be the same since we are collecting license information
		// from the project data itself (the installed package files).
		packages = append(packages, model.Package{
			Description:      spdxhelpers.Description(p),
			DownloadLocation: spdxhelpers.DownloadLocation(p),
			ExternalRefs:     spdxhelpers.ExternalRefs(p),
			FilesAnalyzed:    false,
			HasFiles:         fileIDs,
			Homepage:         spdxhelpers.Homepage(p),
			LicenseDeclared:  license, // The Declared License is what the authors of a project believe govern the package
			Originator:       spdxhelpers.Originator(p),
			SourceInfo:       spdxhelpers.SourceInfo(p),
			VersionInfo:      p.Version,
			Item: model.Item{
				LicenseConcluded: license, // The Concluded License field is the license the SPDX file creator believes governs the package
				Element: model.Element{
					SPDXID: packageSpdxID,
					Name:   p.Name,
				},
			},
		})
	}

	return packages, files, relationships
}

func cleanSPDXName(name string) string {
	// remove # according to specification
	name = strings.ReplaceAll(name, "#", "-")

	// remove : for url construction
	name = strings.ReplaceAll(name, ":", "-")

	// clean relative pathing
	return path.Clean(name)
}
