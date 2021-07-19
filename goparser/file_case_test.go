// ssss
// bvbbb
package goparser

import "context"

const (

	A S = 1 //A sss
)

var (
	bb string = "bb" //bb xx
)

type S int
type SP *int

type SS struct {
	Id string `json:"id,omitempty"`
	Name *SS
}

type ff func(ctx context.Context, ss *SS, bb SS) (err error)

type tt interface {
	A()
}