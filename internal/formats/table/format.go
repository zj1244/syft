package table

import "github.com/zj1244/syft/syft/format"

func Format() format.Format {
	return format.NewFormat(
		format.TableOption,
		encoder,
		nil,
		nil,
	)
}
