##Top-Ten-Words

####What is it?

A REST API built with go standard library that accepts large size of text(10 million chars), concurrently processes it and give json response of it's top ten words with frequency.

#### How to use it?

* run the server from root directory with `go run main.go`.
* send a post request to `http://localhost:8080/topten` with `text` as a body.


#### Directory Structure

* `internal` directory has `dto` package consists of text_request from client and top_words_response  and `lib` package has common files like constants, structured error library.
* `pkg` directory has `service` consists of core top_words service and it's helpers.
* `cmd` directory has `main.go` which is the entry point of the application.

#### Few gotchas:

* Text will be processed concurrently utilizing maximum number of cpu cores.
* If GoMaxProcs is 4, then 4 concurrent goroutines will divide text by 4 chunks and processes it all together.
* Initially, concurrency concepts used so far, sync wait group,mutex, goroutines.