module github.com/clubrizer/services/hello

go 1.19

//replace github.com/clubrizer/server/pkg => ./../../pkg

require (
	github.com/clubrizer/server/pkg v0.1.2
	github.com/go-chi/chi/v5 v5.0.7
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
)
