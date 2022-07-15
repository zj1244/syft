package php

import (
	"os"
	"testing"

	"github.com/zj1244/syft/syft/pkg"
	"github.com/go-test/deep"
)

func TestParseComposerFileLock(t *testing.T) {
	expected := []pkg.Package{
		{
			Name:     "adoy/fastcgi-client",
			Version:  "1.0.2",
			Language: pkg.PHP,
			Type:     pkg.PhpComposerPkg,
		},
		{
			Name:     "alcaeus/mongo-php-adapter",
			Version:  "1.1.11",
			Language: pkg.PHP,
			Type:     pkg.PhpComposerPkg,
		},
	}
	fixture, err := os.Open("test-fixtures/composer.lock")
	if err != nil {
		t.Fatalf("failed to open fixture: %+v", err)
	}

	actual, err := parseComposerLock(fixture.Name(), fixture)
	if err != nil {
		t.Fatalf("failed to parse requirements: %+v", err)
	}
	differences := deep.Equal(expected, actual)
	if differences != nil {
		t.Errorf("returned package list differed from expectation: %+v", differences)
	}

}
