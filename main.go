package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	url := flag.String("url", "", "remote proxy url")
	regions := flag.String("regions", "", "regions, separated with comma")
	path := flag.String("path", "/Users/yee/Library/Mobile Documents/iCloud~run~surge/Documents/", "local file path")
	flag.Parse()

	if "" == *url {
		log.Fatal("url param is invalid")
	} else if "" == *regions {
		log.Fatal("regions param is invalid")
	}

	log.Printf("url: %s\n", *url)
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*20)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatalf("httpClient.Do() failed with:\n'%s'\n", err)
	}
	defer resp.Body.Close()

	res, err := groupProxies(resp.Body, strings.Split(*regions, ","))

	// log.Printf("%+v", res)

	if err = writeProxyFiles(*path, res); err != nil {
		log.Panic(err)
	}
}

func groupProxies(r io.Reader, g []string) (map[string][]string, error) {
	res := make(map[string][]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		for _, key := range g {
			// log.Printf("%+v", res[key])
			if strings.Contains(l, key) {
				log.Printf("current key:%s, matched line: %s\n", key, l)
				res[key] = append(res[key], l)
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func writeProxyFiles(path string, m map[string][]string) error {
	for k, v := range m {
		filename := path + k + ".list"
		if _, err := os.Stat(filename); err == nil {
			if err = os.Remove(filename); err != nil {
				log.Fatal(err)
			}
		}

		if len(v) == 0 {
			log.Printf("%s doesn't contain any server.\n", k)
			continue
		}

		log.Printf("file:%s\n", filename)
		f, err := os.Create(filename)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		for _, l := range v {
			w.WriteString(l)
			w.WriteString("\n")
		}
		w.Flush()
	}

	return nil
}