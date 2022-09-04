module github.com/clubrizer/server/env

go 1.19

//replace github.com/clubrizer/log => ../log

require (
	github.com/clubrizer/server/log v0.1.0
	github.com/joho/godotenv v1.4.0
)

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
)
