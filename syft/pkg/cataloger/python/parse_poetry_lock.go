package python

import (
	"fmt"
	"io"

	"github.com/pelletier/go-toml"
	"github.com/zj1244/syft/syft/pkg"
	"github.com/zj1244/syft/syft/pkg/cataloger/common"
)

// integrity check
var _ common.ParserFn = parsePoetryLock

// parsePoetryLock is a parser function for poetry.lock contents, returning all python packages discovered.
func parsePoetryLock(_ string, reader io.Reader) ([]pkg.Package, error) {
	tree, err := toml.LoadReader(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to load poetry.lock for parsing: %v", err)
	}

	metadata := PoetryMetadata{}
	err = tree.Unmarshal(&metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to parse poetry.lock: %v", err)
	}

	return metadata.Pkgs(), nil
}
