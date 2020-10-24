# Multithreaded Webcrawler

An implementation of a multithreaded web crawler in Go. To add extra challenge, the `getURLs` function adds latency to test how concurrency speeds this up. Also note that read and write functions are thread-safe with the inclusion of a reader writer mutex. 

## To Run
```
go build main.go
./main
```
