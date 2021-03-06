package deb

import (
	"bufio"
	"os"
	"testing"

	"github.com/zj1244/syft/syft/file"

	"github.com/go-test/deep"
	"github.com/zj1244/syft/syft/pkg"
)

func compareEntries(t *testing.T, left, right pkg.DpkgMetadata) {
	t.Helper()
	if diff := deep.Equal(left, right); diff != nil {
		t.Error(diff)
	}
}

func TestSinglePackage(t *testing.T) {
	tests := []struct {
		name     string
		expected pkg.DpkgMetadata
	}{
		{
			name: "Test Single Package",
			expected: pkg.DpkgMetadata{
				Package:       "apt",
				Source:        "apt-dev",
				Version:       "1.8.2",
				Architecture:  "amd64",
				InstalledSize: 4064,
				Maintainer:    "APT Development Team <deity@lists.debian.org>",
				Files: []pkg.DpkgFileRecord{
					{
						Path: "/etc/apt/apt.conf.d/01autoremove",
						Digest: &file.Digest{
							Algorithm: "md5",
							Value:     "76120d358bc9037bb6358e737b3050b5",
						},
						IsConfigFile: true,
					},
					{
						Path: "/etc/cron.daily/apt-compat",
						Digest: &file.Digest{
							Algorithm: "md5",
							Value:     "49e9b2cfa17849700d4db735d04244f3",
						},
						IsConfigFile: true,
					},
					{
						Path: "/etc/kernel/postinst.d/apt-auto-removal",
						Digest: &file.Digest{
							Algorithm: "md5",
							Value:     "4ad976a68f045517cf4696cec7b8aa3a",
						},
						IsConfigFile: true,
					},
					{
						Path: "/etc/logrotate.d/apt",
						Digest: &file.Digest{
							Algorithm: "md5",
							Value:     "179f2ed4f85cbaca12fa3d69c2a4a1c3",
						},
						IsConfigFile: true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.Open("test-fixtures/status/single")
			if err != nil {
				t.Fatal("Unable to read test_fixtures/single: ", err)
			}
			defer func() {
				err := file.Close()
				if err != nil {
					t.Fatal("closing file failed:", err)
				}
			}()

			reader := bufio.NewReader(file)

			entry, err := parseDpkgStatusEntry(reader)
			if err != nil {
				t.Fatal("Unable to read file contents: ", err)
			}

			compareEntries(t, entry, test.expected)
		})
	}
}

func TestMultiplePackages(t *testing.T) {
	tests := []struct {
		name     string
		expected []pkg.DpkgMetadata
	}{
		{
			name: "Test Multiple Package",
			expected: []pkg.DpkgMetadata{
				{
					Package:       "tzdata",
					Version:       "2020a-0+deb10u1",
					Source:        "tzdata-dev",
					Architecture:  "all",
					InstalledSize: 3036,
					Maintainer:    "GNU Libc Maintainers <debian-glibc@lists.debian.org>",
					Files:         []pkg.DpkgFileRecord{},
				},
				{
					Package:       "util-linux",
					Version:       "2.33.1-0.1",
					Architecture:  "amd64",
					InstalledSize: 4327,
					Maintainer:    "LaMont Jones <lamont@debian.org>",
					Files: []pkg.DpkgFileRecord{
						{
							Path: "/etc/default/hwclock",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "3916544450533eca69131f894db0ca12",
							},
							IsConfigFile: true,
						},
						{
							Path: "/etc/init.d/hwclock.sh",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "1ca5c0743fa797ffa364db95bb8d8d8e",
							},
							IsConfigFile: true,
						},
						{
							Path: "/etc/pam.d/runuser",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "b8b44b045259525e0fae9e38fdb2aeeb",
							},
							IsConfigFile: true,
						},
						{
							Path: "/etc/pam.d/runuser-l",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "2106ea05877e8913f34b2c77fa02be45",
							},
							IsConfigFile: true,
						},
						{
							Path: "/etc/pam.d/su",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "ce6dcfda3b190a27a455bb38a45ff34a",
							},
							IsConfigFile: true,
						},
						{
							Path: "/etc/pam.d/su-l",
							Digest: &file.Digest{
								Algorithm: "md5",
								Value:     "756fef5687fecc0d986e5951427b0c4f",
							},
							IsConfigFile: true,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.Open("test-fixtures/status/multiple")
			if err != nil {
				t.Fatal("Unable to read: ", err)
			}
			defer func() {
				err := file.Close()
				if err != nil {
					t.Fatal("closing file failed:", err)
				}
			}()

			pkgs, err := parseDpkgStatus(file)
			if err != nil {
				t.Fatal("Unable to read file contents: ", err)
			}

			if len(pkgs) != 2 {
				t.Fatalf("unexpected number of entries: %d", len(pkgs))
			}

			for idx, entry := range pkgs {
				compareEntries(t, entry.Metadata.(pkg.DpkgMetadata), test.expected[idx])
			}

		})
	}
}

func TestSourceVersionExtract(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "name and version",
			input:    "test (1.2.3)",
			expected: []string{"test", "1.2.3"},
		},
		{
			name:     "only name",
			input:    "test",
			expected: []string{"test", ""},
		},
		{
			name:     "empty",
			input:    "",
			expected: []string{"", ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, version := extractSourceVersion(test.input)

			if name != test.expected[0] {
				t.Errorf("mismatch name for %q : %q!=%q", test.input, name, test.expected[0])
			}

			if version != test.expected[1] {
				t.Errorf("mismatch version for %q : %q!=%q", test.input, version, test.expected[1])
			}

		})
	}
}
