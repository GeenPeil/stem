
## rutte
The backend service for the GeenPeil platform.

### Quickstart
- Install [Go](https://golang.org), follow [these instructions](https://golang.org/doc/install).
- Make sure you have checked out this repository in your GOPATH.
- Install modd: `go get github.com/cortesi/modd/cmd/modd`
- Change your working directory to this folder (`cd $GOPATH/src/github.com/GeertPeil/rutte`)
- Run go get to fetch dependencies (`go get -u -a`)
- Run modd in this directory (`rutte`): `modd`

You should now have `rutte` running and listening for API calls. `modd` will rebuild and restart `rutte` whenever a `.go` file is modified.
