package text

import "github.com/zj1244/syft/syft/format"

func Format() format.Format {
	return format.NewFormat(
		format.TextOption,
		encoder,
		nil,
		nil,
	)
}
