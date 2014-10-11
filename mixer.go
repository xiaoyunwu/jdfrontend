package main

import (
    "net/http"
    "net/url"
    "io"
    "log"
)

func main() {
    tr := &http.Transport{
        DisableCompression: true,
        DisableKeepAlives: false,
        MaxIdleConnsPerHost: 48,
    }

    client := &http.Client{Transport: tr}

    http.HandleFunc("/jds", func(w http.ResponseWriter, r *http.Request) {
        proxy(w, r, client)
    }) 

    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// For now, we just need to copy the parameters over and copy the result back.
// Late, we can add cache support, and also a/b testing support, and ads mixing
// support. We also need to add tracking support.
func proxy(w http.ResponseWriter, r *http.Request, client *http.Client) {
    // Pointing to the real backend.
    var lurl *url.URL
    lurl, err := url.Parse("http://localhost:9200")
    lurl.Path += "/jidian/_jds"
    parameters := url.Values{}
    
    // Now we forward all the parameter from the client.
    r.ParseForm()
    for k, vs := range r.Form {
        for _, v := range vs {
            if k
            parameters.Add(k, v)
        }
    }
        
    lurl.RawQuery = parameters.Encode()

    // Fetch the result from the real backend
    req, err := http.NewRequest("GET", lurl.String(), nil)
    if err != nil {
        log.Fatal(err)
        return
    }

    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
        w.WriteHeader(resp.StatusCode)
        return
    }       

    // Copy the content to client.
    copyHeaders(w.Header(), resp.Header)
    w.WriteHeader(resp.StatusCode)
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        log.Fatal(err)
        return
    }
}  

func copyHeaders(dst, src http.Header) {
    for k, _ := range dst {
        dst.Del(k)
    }
    for k, vs := range src {
        for _, v := range vs {
            dst.Add(k, v)
        }
    }
}