RUN Server

```
$ go run -race hashserver.go handlers.go
2021/10/16 16:16:56 Worker Running
2021/10/16 16:17:10 Request duration: 285.117µs &{POST /hash HTTP/1.1 1 1 map[Accept:[*/*] Content-Length:[21] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.64.0]] 0xc000056300 <nil> 21 [] false localhost:8080 map[password:[angryMonkey1]] map[password:[angryMonkey1]] <nil> map[] 127.0.0.1:47796 /hash <nil> <nil> <nil> 0xc000056340}
2021/10/16 16:17:10 Run job {1 angryMonkey1}
2021/10/16 16:17:10 Request duration: 245.432µs &{POST /hash HTTP/1.1 1 1 map[Accept:[*/*] Content-Length:[21] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.64.0]] 0xc000056480 <nil> 21 [] false localhost:8080 map[password:[angryMonkey2]] map[password:[angryMonkey2]] <nil> map[] 127.0.0.1:47798 /hash <nil> <nil> <nil> 0xc0000564c0}
2021/10/16 16:17:10 Run job {2 angryMonkey2}
...
2021/10/16 16:17:15 Hash duration: 79.697µs angryMonkey1 ea7KzzR5nQ70wpEAzVNokIf0LGbkd9/MhBSIelTDfA9NObui2QIgceeGtmWfXbMAWlm+c2OwecXkyGq2UYrZsA==
2021/10/16 16:17:15 {1 angryMonkey1} ea7KzzR5nQ70wpEAzVNokIf0LGbkd9/MhBSIelTDfA9NObui2QIgceeGtmWfXbMAWlm+c2OwecXkyGq2UYrZsA==
2021/10/16 16:17:15 Hash duration: 59.846µs angryMonkey2 k9w4qLUoXzJEUp6TXTL59Vhjq1g600F0Va9v/VLkNeegC7Oro7kh/AIMU20+RlnG4fBDdfmv9qY4NHc5rF7YTw==
2021/10/16 16:17:15 {2 angryMonkey2} k9w4qLUoXzJEUp6TXTL59Vhjq1g600F0Va9v/VLkNeegC7Oro7kh/AIMU20+RlnG4fBDdfmv9qY4NHc5rF7YTw==
2021/10/16 16:17:15 Hash duration: 226.523µs angryMonkey3 +OKDXmFdy5f2WBlJlxckDJpZPho0vEhMs6h5luF5fCOKqnFTluWQdDU2eMoPfkDNbh+tI7ANiSjFIwFD8wDcbA==
...
```

Client:

## POST /hash
COMPUTE HASH
```
$ for i in {1..100}; do curl --data "password=angryMonkey$i" localhost:8080/hash; done
counter:1
counter:2
...
```

## GET /hash/\<reqid\>
```
$ curl localhost:8080/hash/1
+OKDXmFdy5f2WBlJlxckDJpZPho0vEhMs6h5luF5fCOKqnFTluWQdDU2eMoPfkDNbh+tI7ANiSjFIwFD8wDcbA==
```

## GET /stats
```
$ while :;do curl localhost:8080/stats; done
{"totalRequests":207,"avgDuration":"249.111µs"}
{"totalRequests":208,"avgDuration":"248.575µs"}
```

# Concurrent Load
```
sh ab.sh
```
TODO:
1. graceful shutdown, https://pkg.go.dev/net/http#Server.Shutdown
2. use context package, pass the cancellation context through to the goroutines.
3. use error values, better error handling
4. better http status codes
5. unit tests
6. log stats, eg: number of goroutines, mem, cpu etc.
