package integration

import (
	"testing"
)

func TestRegressionJavaNoMainPackage(t *testing.T) { // Regression: https://github.com/zj1244/syft/issues/252
	catalogFixtureImage(t, "image-java-no-main-package")
}
