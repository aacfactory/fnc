# FNC

fns codes generator

## Install
### Download
See [releases](https://github.com/aacfactory/fnc/releases)
### Build from source
```bash
go install github.com/aacfactory/fnc@latest
```
## Usage
### Create project
```bash
cd {your project dir}
fnc create -p {project mod path} .
```
### Generate codes
mark `go:generate`
```go
// main
// go:generate fnc codes .
func main() {
	
}
```
run generate
```bash
go generate
```


