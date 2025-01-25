package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mirkobrombin/copr2url/structs"
	"github.com/mirkobrombin/copr2url/utils"
	"golang.org/x/net/html"
	"gopkg.in/ini.v1"
)

func main() {
	// Parsing arguments
	iniFile := "repos.ini"
	fedoraTarget := "fedora-rawhide-x86_64"
	doJSON := false

	var positional []string
	for _, arg := range os.Args[1:] {
		if arg == "--json" {
			doJSON = true
		} else if strings.HasPrefix(arg, "--") {
			continue // simply ignoring unknown flags
		} else {
			positional = append(positional, arg)
		}
	}
	if len(positional) > 0 {
		iniFile = positional[0]
	}
	if len(positional) > 1 {
		fedoraTarget = positional[1]
	}

	// Loading repos from .ini
	cfg, err := ini.Load(iniFile)
	if err != nil {
		log.Fatalf("Cannot load ini: %v", err)
	}

	var repos []structs.Repo
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}
		repos = append(repos, structs.Repo{
			Owner:   section.Key("owner").String(),
			Project: section.Key("project").String(),
			Package: section.Key("package").String(),
		})
	}

	// Collecting final links
	var links []string
	for _, r := range repos {
		buildID := getLatestSuccessBuild(r.Owner, r.Project, r.Package)
		if buildID == 0 {
			log.Printf("No success build for %s/%s/%s\n", r.Owner, r.Project, r.Package)
			continue
		}
		rpmLink := getRpmFromListing(r.Owner, r.Project, r.Package, fedoraTarget, buildID)
		if rpmLink == "" {
			log.Printf("No matching RPM for %s (build %d)\n", r.Package, buildID)
			continue
		}
		links = append(links, rpmLink)
	}

	// Output
	if doJSON {
		out, _ := json.MarshalIndent(links, "", "  ")
		fmt.Println(string(out))
	} else {
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

// getLatestSuccessBuild returns the most recent successful build ID for a given package.
func getLatestSuccessBuild(owner, project, pkg string) int {
	url := fmt.Sprintf("https://copr.fedorainfracloud.org/api_3/build/list?ownername=%s&projectname=%s", owner, project)
	body, err := utils.FetchBody(url)
	if err != nil {
		return 0
	}
	var resp structs.BuildListResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0
	}
	var ids []int
	for _, item := range resp.Items {
		if item.State == "succeeded" && strings.EqualFold(item.SourcePackage.Name, pkg) {
			ids = append(ids, item.ID)
		}
	}
	if len(ids) == 0 {
		return 0
	}
	sort.Ints(ids)
	return ids[len(ids)-1]
}

// getRpmFromListing parses an HTML directory listing to find a suitable .rpm file.
func getRpmFromListing(owner, project, pkg, fedoraTarget string, buildID int) string {
	// Extracting desired arch from fedoraTarget, e.g. "fedora-41-x86_64" -> "x86_64"
	var archWanted string
	if strings.HasSuffix(fedoraTarget, "x86_64") {
		archWanted = "x86_64"
	} else if strings.HasSuffix(fedoraTarget, "aarch64") {
		archWanted = "aarch64"
	}

	// Building the directory URL
	url := fmt.Sprintf("https://download.copr.fedorainfracloud.org/results/%s/%s/%s/0%d-%s/",
		owner, project, fedoraTarget, buildID, pkg)
	body, err := utils.FetchBody(url)
	if err != nil {
		return ""
	}
	var rpmLinks []string

	// Parsing HTML looking for <a href="...rpm"> entries
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return ""
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			h := utils.Attr(n, "href")
			if strings.HasSuffix(h, ".rpm") {
				if strings.Contains(h, "debug") || strings.Contains(h, "src.rpm") {
					return
				}
				rpmLinks = append(rpmLinks, h)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	if len(rpmLinks) == 0 {
		return ""
	}

	// Filtering by arch or noarch
	var archMatches, noarchMatches []string
	for _, link := range rpmLinks {
		full := url + link
		if strings.Contains(link, ".noarch.rpm") {
			noarchMatches = append(noarchMatches, full)
		} else if archWanted != "" && strings.Contains(link, "."+archWanted+".rpm") {
			archMatches = append(archMatches, full)
		}
	}
	sort.Strings(archMatches)
	sort.Strings(noarchMatches)

	if len(archMatches) > 0 {
		return archMatches[0]
	}
	if len(noarchMatches) > 0 {
		return noarchMatches[0]
	}
	// If archWanted is empty, just return the first .rpm
	if archWanted == "" {
		sort.Strings(rpmLinks)
		return url + rpmLinks[0]
	}
	return ""
}
