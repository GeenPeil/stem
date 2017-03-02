
## Development on windows

This guide explains how to get up and running on windows.


### Install base tools
First, we need some base tools. Download and install:
- [GitHub Desktop](https://desktop.github.com/)
- [Golang 64-bit MSI installer](https://golang.org/dl/)
- [NodeJS 64-bit MSI installer](https://nodejs.org/en/download/)
- [PostgreSQL 9.6.x](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads#windows) (After the main installation has finished, you don't need to launch Stack Driver for additional tools)
- [Make](http://gnuwin32.sourceforge.net/packages/make.htm) ([direct download installer](http://gnuwin32.sourceforge.net/downlinks/make.php))

We'll also need an editor/IDE. Personally I prefer [Atom](https://atom.io) with a range of plugins such as go-plus and typescript. But I also hear good things about Visual Studio and VSCode.

### Clone repository
Open GitHub Desktop and log in with your github account. Open Git Shell, this was installed with GitHub Desktop. 

Run the following commands:
```
cd \
mkdir Gopath\src\github.com\GeenPeil
cd Gopath\src\github.com\GeenPeil
git clone git@github.com:GeenPeil/stem.git
```

### Add Gopath to windows environment
Close any open shells to avoid annoying mistakes later. Add a new entry to your system environment variables. The name for this variable must be `GOPATH` and the value is `C:\Gopath`. Also modify the `Path` environment variable. Add the following two values at the end of the list:
- `C:\Gopath\bin`
- `C:\Program Files (x86)\GnuWin32\bin`

### Install tools
Open a new Git Shell. We need a **new** shell because we just modified the environment variables. Execute:
```
npm install -g typescript
go get github.com/cortesi/modd/cmd/modd
```

### Run rutte
1) Open Git shell and go the the `rutte` directory
```
cd \Gopath\src\github.com\GeenPeil\stem\rutte
```
2) Now we need to download the dependencies for the go code in the `rutte` directory, and its subdirectories. This can take a while.
```
make dependencies
```
3) Next, build the project
```
make build
```
4) Lastly, run rutte. Note that `rutte.exe` will probably exit now because the database has not been set up yet.
```
.\dist\bin\rutte.exe
```


If you plan to modify Go code, stop `rutte.exe` (`ctrl + c`) and run `modd`. The tool `modd` will watch for changes to `.go` files and automatically rebuilds and restarts `rutte.exe`.
