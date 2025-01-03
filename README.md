# Wurdsearch - Search Engine Project

![demo](./docs/demo.gif)

## Usage

1. Make sure MongoDB is [installed](https://www.mongodb.com/docs/manual/installation/) and running on port 27017.

2. Crawl some websites with `go run crawler/cmd/main.go`. Set up the seed websites and crawl depth in `main.go`

3. Index the crawled data: `go run indexer/cmd/main.go`

4. Start the backend: `go run search/cmd/main.go`

5. Start the frontend: `cd frontend && npm install && npm run dev`

6. Visit [http://localhost:5173](http://localhost:5173/) to use Wurdsearch!

![logo](./frontend/assets/images/web-app-manifest-192x192.png)
