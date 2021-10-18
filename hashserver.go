package main

import (
    "log"
    "time"
    "sync"
    "crypto/sha512"
    "encoding/base64"
    "strconv"
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

