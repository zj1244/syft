package syftjson

import "github.com/zj1244/syft/syft/format"

func Format() format.Format {
	return format.NewFormat(
		format.JSONOption,
		encoder,
		decoder,
		validator,
	)
}
