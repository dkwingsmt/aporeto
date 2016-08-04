package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "flag"
    "strings"
    "unicode"
    "time"
)

const usage string = `gosample [-help|-h]
gosample -urls=<comma-seperated-one-or-more-urls>`

func main() {
    urlsstr := flag.String("urls", "", "Comma separated urls to fetch")
    flag.Parse()
    urls := strings.Split(*urlsstr, ",")
    durationCh := make(chan int64)
    totalTime := int64(0)
    start := time.Now()
    for i, url := range urls {
        go fetch(url, i, durationCh)
    }
    for range urls {
        totalTime += <-durationCh
    }
    fmt.Println("Actual time:", time.Since(start).Nanoseconds())
    fmt.Println("Total individual time:", totalTime)
}

func fetch(url string, i int, duration chan<- int64) {
    start := time.Now()
    // Get url body
    resp, err := http.Get(url)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    // Count word frequency
    // Split a string by all non-alnum chars
    words := strings.FieldsFunc(string(body), func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
    dict := make(map[string]int)
    for _, string := range words {
        _, present := dict[string]
        if present {
            dict[string]++
        } else {
            dict[string] = 1
        }
    }

    // Write result to file
    filename := fmt.Sprintf("url%d.txt", i + 1)
    f, err := os.Create(filename)
    if err != nil {
        return
    }
    defer f.Close()
    fmt.Fprintln(f, "url:", url)
    for word, count := range dict {
        fmt.Fprintf(f, "  %s: %d\n", word, count)
    }
    duration <- time.Since(start).Nanoseconds()
}
