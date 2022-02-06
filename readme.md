## Top Ten Words

#### What is it?

A REST API built with go standard library that accepts large size of text(10 million chars), concurrently processes it and give json response of it's top ten words with frequency.

#### How to use it?

* run the server from root directory with `go run main.go`.
* send a post request to `http://localhost:8080/topten` with `text` as `JSON` in body.


#### Directory Structure

* `internal` directory has
  * `dto` package consists of `text_request` from client, `top_words_response` to client.
  * `lib` package has common files like `constants`, structured `errors` library.
* `pkg` directory has 
  * `service` package consists of core `top_words service`, it's `helpers` and `json_utils`.
* `cmd` directory has 
  * `api` package consists of `app.go` and it's REST `handlers`.

#### Data Flow

Client --(JSON)-> REST Handlers --(DTO)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

#### Few gotchas:

* Text will be processed concurrently utilizing maximum number of cpu cores.
* If GoMaxProcs is 4, then 4 concurrent goroutines(workers) will divide all words by chunks and will be executed concurrently.
* Initially, concurrency concepts used so far, sync wait group,mutex lock unlock, goroutines.
