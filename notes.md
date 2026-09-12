
# Handler

A "handler" is just a function that runs when someone visits a specific URL.

```
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Golang!")
}
```

- `w http.ResponseWriter` is how you write the response back to the browser
- `r *http.Request` holds info about the incoming request (URL, headers, etc)


# Start the Server
```
http.HandleFunc("/", homeHandler)

if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Fatal(err)
}
```
`http.ListenAndServe(":8080", nil)` actually starts the web server, listening on port 8080. nil means "use the default router" (the one we registered routes on above via http.HandleFunc).