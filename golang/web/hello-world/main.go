package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const indexTemplate = `
<html>
	<head>
		<meta charset="utf-8">
		<title>Main Page</title>
	</head>
	<body>
		<h1>Welcome to the simplest web app!</h1>
	</body>
</html>
`

var index = template.Must(template.New("index").Parse(indexTemplate))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	addr := "localhost:8101"
	fmt.Printf("Starting HTTP server at %q", addr)

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
