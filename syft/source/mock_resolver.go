package source

import (
	"fmt"
	"io"
	"os"

	"github.com/anchore/syft/internal/file"
)

var _ Resolver = (*MockResolver)(nil)

// MockResolver implements the Resolver interface and is intended for use *only in test code*.
// It provides an implementation that can resolve local filesystem paths using only a provided discrete list of file
// paths, which are typically paths to test fixtures.
type MockResolver struct {
	Locations []Location
}

// NewMockResolverForPaths creates a new MockResolver, where the only resolvable
// files are those specified by the supplied paths.
func NewMockResolverForPaths(paths ...string) *MockResolver {
	var locations []Location
	for _, p := range paths {
		locations = append(locations, NewLocation(p))
	}

	return &MockResolver{Locations: locations}
}

// HasPath indicates if the given path exists in the underlying source.
func (r MockResolver) HasPath(path string) bool {
	for _, l := range r.Locations {
		if l.RealPath == path {
			return true
		}
	}
	return false
}

// String returns the string representation of the MockResolver.
func (r MockResolver) String() string {
	return fmt.Sprintf("mock:(%s,...)", r.Locations[0].RealPath)
}

// FileContentsByLocation fetches file contents for a single location. If the
// path does not exist, an error is returned.
func (r MockResolver) FileContentsByLocation(location Location) (io.ReadCloser, error) {
	for _, l := range r.Locations {
		if l == location {
			return os.Open(location.RealPath)
		}
	}

	return nil, fmt.Errorf("no file for location: %v", location)
}

// MultipleFileContentsByLocation returns the file contents for all specified Locations.
func (r MockResolver) MultipleFileContentsByLocation(locations []Location) (map[Location]io.ReadCloser, error) {
	results := make(map[Location]io.ReadCloser)
	for _, l := range locations {
		contents, err := r.FileContentsByLocation(l)
		if err != nil {
			return nil, err
		}
		results[l] = contents
	}

	return results, nil
}

// FilesByPath returns all Locations that match the given paths.
func (r MockResolver) FilesByPath(paths ...string) ([]Location, error) {
	var results []Location
	for _, p := range paths {
		for _, location := range r.Locations {
			if p == location.RealPath {
				results = append(results, NewLocation(p))
			}
		}
	}

	return results, nil
}

// FilesByGlob returns all Locations that match the given path glob pattern.
func (r MockResolver) FilesByGlob(patterns ...string) ([]Location, error) {
	var results []Location
	for _, pattern := range patterns {
		for _, location := range r.Locations {
			if file.GlobMatch(pattern, location.RealPath) {
				results = append(results, location)
			}
		}
	}

	return results, nil
}

// RelativeFileByPath returns a single Location for the given path.
func (r MockResolver) RelativeFileByPath(_ Location, path string) *Location {
	paths, err := r.FilesByPath(path)
	if err != nil {
		return nil
	}

	if len(paths) < 1 {
		return nil
	}

	return &paths[0]
}
