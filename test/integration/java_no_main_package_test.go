package integration

import (
	"github.com/anchore/stereoscope/pkg/imagetest"
	"github.com/zj1244/syft/syft"
	"github.com/zj1244/syft/syft/source"
	"testing"
)

func TestJavaNoMainPackage(t *testing.T) { // Regression: https://github.com/zj1244/syft/issues/252
	fixtureImageName := "image-java-no-main-package"
	_, cleanup := imagetest.GetFixtureImage(t, "docker-archive", fixtureImageName)
	tarPath := imagetest.GetFixtureImageTarPath(t, fixtureImageName)
	defer cleanup()

	_, _, _, err := syft.Catalog("docker-archive:"+tarPath, source.SquashedScope)
	if err != nil {
		t.Fatalf("failed to catalog image: %+v", err)
	}
}
