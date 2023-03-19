# FNC

fns codes generator

## Install
### Download
See [releases](https://github.com/aacfactory/fnc/releases)
### Build from source
```bash
go get github.com/aacfactory/fnc
go install github.com/aacfactory/fnc
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

### Mark annotations
Enable service annotations.  
doc.go in fn source file's folder.
```go
// Package samples
// @service samples
// @title samples
// @description samples service
package samples
```
Enable fn annotations.
```go
// query
// @fn query
// @title query
// @description query
func query(ctx context.Context, argument QueryArgument) (result Result, err errors.CodeError) {

	return
}

```
Enable argument and result annotations.
```go
// QueryArgument
// @title title 
// @description description
type QueryArgument struct {
	// Offset
	// @title title 
	// @description description
	Offset int `json:"offset" validate:"required" message:"offset is invalid"`
	// Limit
	// @title title 
	// @description description
	Limit  int `json:"limit" validate:"required" message:"limit is invalid"`
}
```
```go
// Sample
// @title title of sample
// @description >>> description of sample
// > support markdown
// 
// <<<
type Sample struct {
	// Id
	// @title 编号
	// @description 编号
	Id string `json:"id"`
	// Mobile
	// @title 手机号
	// @description 手机号
	Mobile string `json:"mobile"`
	// Name
	// @title 姓名
	// @description 姓名
	Name string `json:"name"`
	// Gender
	// @title 性别
	// @enum M,F,N
	// @description 性别
	Gender string `json:"gender"`
	// Age
	// @title 年龄
	// @description 年龄
	// @enum 1,2,3
	Age int `json:"age"`
	// Avatar
	// @title 头像图片地址
	// @description 头像图片地址
	Avatar string `json:"avatar"`
	// Score
	// @title Score
	// @description Score
	Score float32 `json:"score,omitempty"`
	// DOB
	// @title DOB
	// @description DOB
	DOB json.Date `json:"dob,omitempty"`
	// CreateAT
	// @title CreateAT
	// @description CreateAT
	CreateAT json.Time `json:"createAt,omitempty"`
	// Tokens
	// @title Tokens
	// @description Tokens
	Tokens []string
	// Users
	// @title Users
	// @description Users
	Users []*users.User `json:"users,omitempty"`
	// UserMap
	// @title UserMap
	// @description UserMap
	UserMap map[string]*users.User `json:"userMap,omitempty"`
	// Raw
	// @title Raw
	// @description Raw
	Raw json.RawMessage
}
```

## Note
> Type of Fn Argument must be value object.  
> Builtin, Star Type, Map, Array and Alias Type are not supported.  

> Type of Fn Result must be object type, when use array type, please use ident.
> Builtin, Value Object, Map, Array and Alias Type are not supported.  
