
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


# PATCH nil check

Use `*string` (pointer) fields in a PATCH request struct, not plain `string`.

```
type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
```

- `nil` = client didn't send this field → leave it alone
- non-nil = client sent it (even `""`) → apply it, dereference with `*input.Name`
- plain `string` can't tell "omitted" apart from "sent as empty" — zero value for both is `""`


# Pointers - real use cases

A pointer (`*T`) holds the address of a variable, not a copy. Use one when a function needs to **modify your original variable**, or **avoid copying** something big.

```
json.NewDecoder(r.Body).Decode(&input)
```
Without `&`, Decode gets a copy of `input`. It can fill in the copy, but the copy disappears when Decode returns — your real `input` stays empty. `&input` gives Decode the address, so it writes directly into your variable.

Common real cases:
- `Decode(&input)` — function needs to fill in your struct
- `for i := range list { list[i].Name = x }` — mutate a slice element in place (looping by value copies, by index doesn't)
- `*http.Request` — big struct, passed by pointer to avoid copying it on every request
- `*string` in `UpdateUserRequest` — see PATCH nil check above, not about mutation, just nil-checking