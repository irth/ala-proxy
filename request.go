package main

import "fmt"

var ErrNotArchivable = fmt.Errorf("not archivable")

type Request interface {
	RepoURL(repo string) string
	ArchiveURL(url string) (string, error)
}

func parseURL(url string) (req Request, err error) {
	req, err = parsePackageURL(url)

	if err == ErrNotAPackage {
		req, err = parseOtherURL(url)
		return
	}

	return
}
