package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "sync"
    "sync/atomic"
    "strconv"
    "crypto/sha512"
    "encoding/base64"
    //"encoding/json"
    "strings"
)

var jobs chan HashIn
var counter uint64
var totalDuration uint64
var m sync.Map
var WAIT_SECONDS = time.Second * 5

type HashIn struct {
    count uint64
    password  string
}

type StatsOut struct {
    ReqCount uint64
    TotalDuration time.Duration
}

func worker(jobs <-chan HashIn) {
    log.Println("Worker Running")
    for job := range jobs {
        log.Println("Run job", job)
        go processJob(job)
    }
}

func processJob(job HashIn){
    time.Sleep(WAIT_SECONDS)
    b64 := base64OfSha512(job.password)
    scount := strconv.FormatUint(job.count, 10)
    m.Store(scount, b64)
    log.Println(job, b64)
}
func base64OfSha512(password string) string{
    startTime := time.Now()
    sha_512 := sha512.New()
    sha_512.Write([]byte(password))
    shab64 := base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
    duration := time.Now().Sub(startTime)
    log.Println("Hash duration:", duration, password, shab64)
    return shab64
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

func computeHashHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        password := r.FormValue("password")
        atomic.AddUint64(&counter, 1)
        hashin := HashIn{
            count: counter,
            password: password,
        }
        jobs <- hashin
        scount := strconv.FormatUint(counter, 10)
        fmt.Fprintln(w, "counter:"+ scount)
    default:
        fmt.Fprintln(w, "Only POST request supported. Requested method:"+ r.Method)
    }
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
    var avgDuration time.Duration
    switch r.Method {
    case "GET":
        if counter == 0{
            fmt.Fprintln(w, "{\"totalRequests\":" + strconv.FormatUint(counter,10) + ",\"avgDuration\":\"" +  avgDuration.String() + "\"}")
        } else {
            avgDuration = time.Duration(totalDuration/counter) * time.Nanosecond
            /*statsout := StatsOut{
                ReqCount: counter,
                TotalDuration: duration,
            }
            fmt.Fprintln(w, json.NewEncoder(w).Encode(statsout))
            */
            fmt.Fprintln(w, "{\"totalRequests\":" + strconv.FormatUint(counter,10) + ",\"avgDuration\":\"" +  avgDuration.String() + "\"}")
        }
    default:
        fmt.Fprintln(w, "Only GET request supported. Requested method:"+ r.Method)
    }
}

func getHashHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        p := strings.Split(r.URL.Path, "/")
        if len(p) != 3 {
            fmt.Fprintln(w, "Missing request number. Provided path:"+ r.URL.Path)
        }
        requestNumber := p[2]
        fmt.Println("Get Hash for request:", requestNumber)
        value, ok:= m.Load(requestNumber)
        if ok {
            fmt.Fprintln(w, value.(string))
        } else {
            fmt.Fprintln(w, "MISSING")
        }
    default:
        fmt.Fprintln(w, "Only GET request supported. Requested method:"+ r.Method)
    }
}


func main() {
    jobs = make(chan HashIn, 100)
    go worker(jobs)

    http.Handle("/hash", timer(http.HandlerFunc(computeHashHandler)))
    http.HandleFunc("/hash/", getHashHandler)
    http.HandleFunc("/stats", statsHandler)
    http.ListenAndServe(":8080", nil)
}
