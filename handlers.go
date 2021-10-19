package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync/atomic"
    "strings"
    "time"
    //"encoding/json"
)
func computeHashHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        password := r.FormValue("password")
        reqNumber := atomic.AddUint64(&counter, 1)
        hashin := HashIn{
            count: reqNumber,
            password: password,
        }
        jobs <- hashin
        scount := strconv.FormatUint(reqNumber, 10)
        fmt.Fprintln(w, "requestNumber:"+ scount)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintln(w, "Only POST request supported. Requested method:" + r.Method)
    }
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
    var avgDuration time.Duration
    switch r.Method {
    case "GET":
        localCounter := atomic.LoadUint64(&counter)
        localTotalDuration := atomic.LoadUint64(&totalDuration)
        if localCounter == 0{
            fmt.Fprintln(w, "{\"totalRequests\":0" + ",\"avgDuration\":\"" +  avgDuration.String() + "\"}")
        } else {
            avgDuration = time.Duration(localTotalDuration/localCounter) * time.Nanosecond
            /*statsout := StatsOut{
                ReqCount: localCounter,
                TotalDuration: avgDuration,
            }
            fmt.Fprintln(w, json.NewEncoder(w).Encode(statsout))
            */
            fmt.Fprintln(w, "{\"totalRequests\":" + strconv.FormatUint(localCounter,10) + ",\"avgDuration\":\"" +  avgDuration.String() + "\"}")
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintln(w, "Only GET request supported. Requested method:" + r.Method)
    }
}

func getHashHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        p := strings.Split(r.URL.Path, "/")[1:]
        if len(p) == 2 && p[1] != "" {
            requestNumber := p[1]
            log.Println("Get Hash for request:", requestNumber)
            value, ok:= m.Load(requestNumber)
            if ok {
                fmt.Fprintln(w, value.(string))
            } else {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprintln(w, "Not Processed Request Number:", requestNumber, ". Check request number or try again later")
            }
        } else {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "No request number in path. Provided path:" + r.URL.Path)
            return
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintln(w, "Only GET request supported. Requested method:" +  r.Method)
    }
}

func timer(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        h.ServeHTTP(w, r)
        duration := time.Now().Sub(startTime)
        atomic.AddUint64(&totalDuration, uint64(duration))
        log.Println("Request duration:", duration, r)
    })
}

func main() {
    jobs = make(chan HashIn, 100)
    go worker(jobs)

    http.Handle("/hash", timer(http.HandlerFunc(computeHashHandler)))
    http.HandleFunc("/hash/", getHashHandler)
    http.HandleFunc("/stats", statsHandler)
    http.ListenAndServe(":8080", nil)
}
