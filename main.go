package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type PacmanArchiveProxy struct {
	ArchiveURL string
	RepoURL    string
}

func tryProxy(w http.ResponseWriter, r *http.Request, url string) (int, error) {
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return 500, fmt.Errorf("error creating request: %w", err)
	}

	for k, v := range r.Header {
		req.Header[k] = v
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return 500, fmt.Errorf("error proxying request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("proxying request failed with status %s", resp.Status)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return 500, fmt.Errorf("error copying response body: %w", err)
	}

	return resp.StatusCode, nil
}

func (p *PacmanArchiveProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request", r.URL, r.Host)

	req, err := parseURL(r.URL.String())
	if err != nil {
		log.Printf("error parsing url %s: %s", r.URL, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Trying upstream:", req.RepoURL(p.RepoURL))
	_, err = tryProxy(w, r, req.RepoURL(p.RepoURL))
	if err == nil {
		return
	}

	log.Println("Upstream failed:", err)

	archive, err := req.ArchiveURL(p.ArchiveURL)
	if err != nil {
		log.Printf("no archive url for %s: %s", r.URL, err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	log.Println("Trying archive:", archive)
	code, err := tryProxy(w, r, archive)
	if err != nil {
		log.Printf("Archive failed: %s", err)
		http.Error(w, err.Error(), code)
	}
}

func main() {
	proxy := PacmanArchiveProxy{
		RepoURL:    "https://mirror.rackspace.com/archlinux/$repo/os/$arch",
		ArchiveURL: "https://archive.archlinux.org",
	}

	log.Println("haiiiii :3")
	http.ListenAndServe(":8080", &proxy)

}
