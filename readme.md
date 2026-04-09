# singerbench
is a go library designed to mass tests vless://* and other shared 
proxy links from internet. it is based on singerbox library which
in turn uses singbox (as in-process-library) under the hood.

# note
build with -tags with_utls at least (see tags for building singbox)
I do like this:
```
go run -tags with_utls cmd/singerbench/main.go add-subscription "subscription-url"
go run -tags with_utls cmd/singerbench/main.go measure
go run -tags with_utls cmd/singerbench/main.go print
```