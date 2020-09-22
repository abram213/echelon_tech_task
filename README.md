## Auth web server

This is auth web server built with authboss.

### How to run?

Run `go run main.go` and open http://localhost:8080 in your browser.
  
or using Docker:
1. Run `docker build -t app_ws .`
2. Run `docker run -p 8080:8080 app_ws`
3. Open http://localhost:8080 in your browser