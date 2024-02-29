package main

import (
	"fmt"
	"regexp"
	"strings"
)

type PackageRequest struct {
	Name      string
	Version   string
	Rel       string
	Extension string

	Filename string
	Filepath string
	Repo     string
	Arch     string
}

var _ Request = (*PackageRequest)(nil)

var ErrNotAPackage = fmt.Errorf("not a package")

func parsePackageURL(url string) (*PackageRequest, error) {
	// check if it's a package, other stuff is not archived in ALA
	isPackage, err := regexp.MatchString(`.*\.pkg\.tar\.\w+`, url)
	if err != nil {
		return nil, err
	}
	if !isPackage {
		return nil, ErrNotAPackage
	}

	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid package url: %s", url)
	}

	//   /core/os/x86_64/package-name-version-rel-arch.extension
	//  0 1    2  3      4 (last)
	repo := parts[1]
	arch := parts[3]
	filename := parts[len(parts)-1]
	filePath := strings.Join(parts[4:], "/")

	// package-name-version-rel-arch.extension
	//  0       1       2       3       4
	// -5      -4      -3      -2      -1
	nameParts := strings.Split(filename, "-")
	if len(nameParts) < 4 {
		return nil, fmt.Errorf("invalid package url: %s", url)
	}
	name := strings.Join(nameParts[:len(nameParts)-3], "-")
	version := nameParts[len(nameParts)-3]
	rel := nameParts[len(nameParts)-2]

	nameParts = strings.SplitN(nameParts[len(nameParts)-1], ".", 2)
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("invalid package url: %s", url)
	}
	extension := nameParts[1]

	return &PackageRequest{
		Name:      name,
		Version:   version,
		Rel:       rel,
		Extension: extension,

		Filename: filename,
		Filepath: filePath,
		Repo:     repo,
		Arch:     arch,
	}, nil
}

func (p *PackageRequest) RepoURL(fmtStr string) string {
	url := strings.ReplaceAll(fmtStr, "$repo", p.Repo)
	url = strings.ReplaceAll(url, "$arch", p.Arch)
	return url + "/" + p.Filepath
}

func (p *PackageRequest) ArchiveURL(url string) (string, error) {
	return fmt.Sprintf("%s/packages/%s/%s/%s", url, p.Name[0:1], p.Name, p.Filename), nil
}
