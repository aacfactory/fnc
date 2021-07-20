// ssss
// bvbbb
package goparser

import (
	"context"
	"time"
)

const (

	A S = 1 //A sss
)

var (
	bb string = "bb" //bb xx
	mm map[string]string = map[string]string{"s":"b"}
)

type S int
type SP *int

type SS struct {
	Id string `json:"id,omitempty" xml:"id"`
	Name *SS
	TT time.Time
}

type ff func(ctx context.Context, ss *SS, bb SS, tt time.Time) (err error)

type tt interface {
	A()
}

type Id []string

type IA map[string]interface{}