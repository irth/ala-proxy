package main

import (
	"fmt"
	"strings"
)

type OtherRequest struct {
	Path string
	Repo string
	Arch string
}

var _ Request = (*OtherRequest)(nil)

func parseOtherURL(url string) (*OtherRequest, error) {
	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid package url: %s", url)
	}
	repo := parts[1]
	arch := parts[3]
	path := strings.Join(parts[4:], "/")

	return &OtherRequest{
		Path: path,
		Repo: repo,
		Arch: arch,
	}, nil
}

func (p *OtherRequest) RepoURL(fmtStr string) string {
	url := strings.ReplaceAll(fmtStr, "$repo", p.Repo)
	url = strings.ReplaceAll(url, "$arch", p.Arch)
	return url + "/" + p.Path
}

func (p *OtherRequest) ArchiveURL(url string) (string, error) {
	return "", ErrNotArchivable
}
