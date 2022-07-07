package source

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/anchore/stereoscope/pkg/image"
)

// Scheme represents the optional prefixed string at the beginning of a user request (e.g. "docker:").
type Scheme string

const (
	// UnknownScheme is the default scheme
	UnknownScheme Scheme = "UnknownScheme"
	// DirectoryScheme indicates the source being cataloged is a directory on the root filesystem
	DirectoryScheme Scheme = "DirectoryScheme"
	// ImageScheme indicates the source being cataloged is a container image
	ImageScheme Scheme = "ImageScheme"
)

func detectScheme(fs afero.Fs, imageDetector sourceDetector, userInput string) (Scheme, string, error) {
	if strings.HasPrefix(userInput, "dir:") {
		// blindly trust the user's scheme
		dirLocation, err := homedir.Expand(strings.TrimPrefix(userInput, "dir:"))
		if err != nil {
			return UnknownScheme, "", fmt.Errorf("unable to expand directory path: %w", err)
		}
		return DirectoryScheme, dirLocation, nil
	}

	// we should attempt to let stereoscope determine what the source is first --just because the source is a valid directory
	// doesn't mean we yet know if it is an OCI layout directory (to be treated as an image) or if it is a generic filesystem directory.
	source, imageSpec, err := imageDetector(userInput)
	if err != nil {
		return UnknownScheme, "", fmt.Errorf("unable to detect the scheme from %q: %w", userInput, err)
	}

	if source == image.UnknownSource {
		dirLocation, err := homedir.Expand(userInput)
		if err != nil {
			return UnknownScheme, "", fmt.Errorf("unable to expand potential directory path: %w", err)
		}

		fileMeta, err := fs.Stat(dirLocation)
		if err != nil {
			return UnknownScheme, "", nil
		}

		if fileMeta.IsDir() {
			return DirectoryScheme, dirLocation, nil
		}
		return UnknownScheme, "", nil
	}

	return ImageScheme, imageSpec, nil
}
