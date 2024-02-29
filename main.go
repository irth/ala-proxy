package main

import (
	"log"
	"net/http"
)

type PacmanArchiveProxy struct {
	ArchiveURL string
	RepoURL    string
}

func (p *PacmanArchiveProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request", r.URL, r.Host)

	req, err := parseURL(r.URL.String())
	if err != nil {
		log.Printf("error parsing url %s: %s", r.URL, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("upstream url", req.RepoURL(p.RepoURL))

	archive, err := req.ArchiveURL(p.ArchiveURL)
	if err != nil {
		log.Printf("no archive url for %s: %s", r.URL, err)
	} else {
		log.Println("archive url", archive)
	}

	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func main() {
	proxy := PacmanArchiveProxy{
		RepoURL:    "https://mirror.rackspace.com/archlinux/$repo/os/$arch",
		ArchiveURL: "https://archive.archlinux.org",
	}

	log.Println("haiiiii :3")
	http.ListenAndServe(":8080", &proxy)

}
