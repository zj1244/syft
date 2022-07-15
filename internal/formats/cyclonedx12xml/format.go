package cyclonedx12xml

import "github.com/zj1244/syft/syft/format"

func Format() format.Format {
	return format.NewFormat(
		format.CycloneDxOption,
		encoder,
		nil,
		nil,
	)
}
