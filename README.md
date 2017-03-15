Yservice
========

A simple proxy to [youtube-dl](https://github.com/rg3/youtube-dl)
that enables downloading of video, audio or both of a valid youtube-dl
url.

It is still a work in progress, most of the parameters are hard coded at this
point, the validation is basic among other simplicities. See **TODO** below
for details.

Requires:
---------

- youtube-dl
- golang (for development)

Getting started:
----------------

First get the code the *Go* way.

`go get github.com/kmwenja/yservice`

Then head over to where the code was cloned

`cd $GOPATH/src/github.com/kmwenja/yservice`

Run it

`go run main.go`

Or build it

`go build`
`./yservice`

How it works
------------

It's a server that starts listening on 8080 for a POST request that contains
(in form-data) the url and the download type (`AUDIO, VIDEO, ALL`). The server then
pushes this download to a channel from which a separate goroutine picks it up and
downloads the url according to the type specified. `VIDEO` will run youtube-dl with no
arguments other than `-c`, `AUDIO` will run youtube-dl with an additional `-x` argument,
and `ALL` will do both. The server prints out on stdout various events including
the queueing, starting and success or failure of the downloads.

TODO:
-----

- [ ] trace download failures
- [ ] enable checking of download status
- [ ] validate url to be downloaded as valid youtube-dl
- [ ] choose filetype of download
- [ ] configure the service's parameters eg download folder, server port
- [ ] start explicit builds and releases
