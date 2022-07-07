package integration

import (
	"testing"

	"github.com/anchore/stereoscope/pkg/imagetest"
	"github.com/go-test/deep"
	"github.com/zj1244/syft/syft"
	"github.com/zj1244/syft/syft/distro"
	"github.com/zj1244/syft/syft/source"
)

func TestDistroImage(t *testing.T) {
	fixtureImageName := "image-distro-id"
	_, cleanup := imagetest.GetFixtureImage(t, "docker-archive", fixtureImageName)
	tarPath := imagetest.GetFixtureImageTarPath(t, fixtureImageName)
	defer cleanup()

	_, _, actualDistro, err := syft.Catalog("docker-archive:"+tarPath, source.SquashedScope)
	if err != nil {
		t.Fatalf("failed to catalog image: %+v", err)
	}

	expected, err := distro.NewDistro(distro.Busybox, "1.31.1", "")
	if err != nil {
		t.Fatalf("could not create distro: %+v", err)
	}

	for _, d := range deep.Equal(actualDistro, &expected) {
		t.Errorf("found distro difference: %+v", d)
	}

}
