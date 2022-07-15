package integration

import (
	"testing"

	"github.com/zj1244/syft/syft/pkg"
)

func TestRegression212ApkBufferSize(t *testing.T) {
	// This is a regression test for issue #212 (https://github.com/zj1244/syft/issues/212) in which the apk db could
	// not be processed due to a scanner buffer that was too small
	catalog, _, _ := catalogFixtureImage(t, "image-large-apk-data")

	expectedPkgs := 58
	actualPkgs := 0
	for range catalog.Enumerate(pkg.ApkPkg) {
		actualPkgs += 1
	}

	if actualPkgs != expectedPkgs {
		t.Errorf("unexpected number of APK packages: %d != %d", expectedPkgs, actualPkgs)
	}
}
