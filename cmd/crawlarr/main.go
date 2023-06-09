package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Debug		bool	`json:"debug"`
	BaseUrl		string	`json:"baseUrl"`
	MatchType	string	`json:"matchType"`
	DepthLimit 	int		`json:"depthLimit"`
	Delay		int		`json:"delay"`
}

var (
    links = map[string]bool{}
	processed = map[string]bool{}
    mutex = &sync.Mutex{}
	wg sync.WaitGroup
	// Create a channel that emits a value every config.Delay
	rate <-chan time.Time
)

func addLink(link string) {
    mutex.Lock()
    defer mutex.Unlock()
    if _, ok := links[link]; !ok {
        links[link] = true
    }
}

func checkLink(link string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := links[link]
	return ok
}

func addProcessed(link string) {
    mutex.Lock()
    defer mutex.Unlock()
    if _, ok := processed[link]; !ok {
        processed[link] = true
    }
}

func checkProcessed(link string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := processed[link]
	return ok
}

func processBody(body []byte, currentUrl string, matchType string, re *regexp.Regexp, baseUrl *url.URL) []string {
    matches := re.FindAllSubmatch(body, -1)
    hits := []string{}

    for _, _submatch := range matches {
        submatch := string(_submatch[1])

		unescaped, _ := url.QueryUnescape(submatch)

        u, err := url.Parse(unescaped)
        if err != nil {
            continue
        }

		curr, err := url.Parse(currentUrl)
		if err != nil {
            continue
        }

		if !u.IsAbs() {
            u = curr.ResolveReference(u)
        }

		switch matchType {
		case "SAME_BASE":
			if !strings.HasPrefix(u.String(), baseUrl.String()) {
				continue
			}
		case "SAME_HOST":
			if (u.Hostname() != baseUrl.Hostname()) {
				continue
			}
		case "SAME_ORIGIN":
			uOrigin := u.Hostname()
			baseOrigin := baseUrl.Hostname()

			uParts := strings.Split(u.Host, ".")
			if len(uParts) > 2 {
				uOrigin = strings.Join(uParts[1:], ".")
			}

			baseParts := strings.Split(baseUrl.Host, ".")
			if len(baseParts) > 2 {
				baseOrigin = strings.Join(baseParts[1:], ".")
			}

			if (uOrigin != baseOrigin) {
				continue
			}

		case "DANGEROUS_NO_MATCH_TYPE_ONLY_ENABLE_IF_YOU_KNOW_WHAT_YOURE_DOING":
			// no filtering, can go wrong
		default:
			if !strings.HasPrefix(u.String(), baseUrl.String()) {
				continue
			}
		}

        // Remove any fragment or query parameters from URL
        u.Fragment = ""
        u.RawQuery = ""

        uStr := u.String()

        if (uStr != currentUrl) {
            hits = append(hits, uStr)
        }
    }

    return hits
}

func processLink(current string, depth int, config *Config, re *regexp.Regexp) {
	defer wg.Done() // decrement WaitGroup counter when done

	current, err := url.QueryUnescape(current)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot unescape current URL: %w", err))
		return
	}
	
	if (depth > 0) {
		<-rate // Wait for the rate limit
	}

	resp, err := http.Get(current)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot HTTP GET base URL: %w", err))
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot read HTTP response body: %w", err))
		return
	}

	base, err := url.Parse(config.BaseUrl)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot read HTTP response body: %w", err))
		return
	}

	hits := []string{}
	if ok := checkProcessed(current); !ok {
		hits = processBody(body, current, config.MatchType, re, base)

		addProcessed(current)

		if (config.Debug) {
			fmt.Printf("Found %d new links on %s\nCurrent depth: %d\n", len(hits), current, depth)
		}
	} else {
		if (config.Debug) {
			fmt.Printf("Already visited %s, skipping...\n", current)
		}
		return
	}

	for _, hit := range hits {
		ok := checkLink(hit)
		if (!ok && depth < config.DepthLimit) {
			addLink(hit)

			wg.Add(1) // increment WaitGroup counter

			go processLink(hit, depth + 1, config, re)
		}
	}
}

func main() {
	configFile, err := os.Open("./config.json")
    if err != nil {
        configFile, err = os.Open("../../config.json")
		if err != nil {
			panic(fmt.Errorf("cannot find/open config.json file: %w", err))
		}
    }
    defer configFile.Close()

    decoder := json.NewDecoder(configFile)
    config := Config{}
    err = decoder.Decode(&config)
    if err != nil {
        panic(fmt.Errorf("cannot parse config.json file: %w", err))
    }

	if (!strings.HasSuffix(config.BaseUrl, "/")) {
		config.BaseUrl = config.BaseUrl + "/"
	}

	rate = time.Tick(time.Duration(config.Delay) * time.Millisecond)

	re := regexp.MustCompile(`<a[\w\s="]*href\s*=\s*(?:\"|')([^\"';]*)(?:\"|')`)

	fmt.Println("Crawlarr started!")

	wg.Add(1) // increment WaitGroup counter

	processLink(config.BaseUrl, 0, &config, re)

	wg.Wait() // wait for all goroutines to complete

	// 0644 = rw-,r--,r--
	f, err := os.OpenFile("links_test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("cannot open links.txt: %w", err))
	} else {
		for link, _ := range links {
			_, err = f.WriteString(link + "\n")
			if err != nil {
				panic(fmt.Errorf("cannot write link to links.txt: %w", err))
			}
		}
	}
	defer f.Close()
	
	fmt.Printf("Crawlarr stopped!\nTotal: %d links\nMax depth: %d\n", len(links), config.DepthLimit)
}
