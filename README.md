# web-crawler
A simple web crawler designed to extract text data to train llms

---

## Project design (v0)

The current design choices where based on favoriting simplicity and getting something that works

- The main components in the app are a crawler, fetcher, parser, and persister
- The app uses a master-workers pattern
- The master orchestrate the workers, e.g. create a worker for each url and handle concurrent access to the queue and the visited set of urls.
- The workers (crawlers) passes the url assigned to them through the crawling pipeline
- The fetcher/parser implementation uses the standard library packages
- The persister implementation uses sqlite to store the extract data
- The queue implementation is a simple in memory queue

---

## Running the application

 ### 1- Clone the project
 ```bash
 git clone https://github.com/mehddii/web-crawler.git
 ```

 ### 2- Build the docker image
```bash
docker image build -t crawler .
 ```

 ### 3- Run the container and mount a local folder for SQLite persistence
```bash
docker run -v "$(pwd)/data:/usr/src/app/data" crawler
```

---

## TBD

- [] Write tests
- [] Change the current orchestration model (e.g. each url processed on its own goroutine)
- [] Switch to sqs
- [] Use a database optimized for object sotrage like s3
- [] Support politeness and rate limiting
- [] Mesure performance
