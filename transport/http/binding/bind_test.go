package binding

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	cnferror "github.com/chnxq/xkitmod/errors"
)

type (
	TestBind struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	TestBind2 struct {
		Age int `json:"age"`
	}
)

func TestBindQuery(t *testing.T) {
	type args struct {
		vars   url.Values
		target any
	}

	tests := []struct {
		name string
		args args
		err  error
		want any
	}{
		{
			name: "test",
			args: args{
				vars:   map[string][]string{"name": {"cnf"}, "url": {"https://go-cnf.dev/"}},
				target: &TestBind{},
			},
			err:  nil,
			want: &TestBind{"cnf", "https://go-cnf.dev/"},
		},
		{
			name: "test1",
			args: args{
				vars:   map[string][]string{"age": {"cnf"}, "url": {"https://go-cnf.dev/"}},
				target: &TestBind2{},
			},
			err: cnferror.BadRequest("CODEC", "Field Namespace:age ERROR:Invalid Integer Value 'cnf' Type 'int' Namespace 'age'"),
		},
		{
			name: "test2",
			args: args{
				vars:   map[string][]string{"age": {"1"}, "url": {"https://go-cnf.dev/"}},
				target: &TestBind2{},
			},
			err:  nil,
			want: &TestBind2{Age: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindQuery(tt.args.vars, tt.args.target)
			if !cnferror.Is(err, tt.err) {
				t.Fatalf("BindQuery() error = %v, err %v", err, tt.err)
			}
			if err == nil && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindQuery() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func TestBindForm(t *testing.T) {
	type args struct {
		req    *http.Request
		target any
	}

	tests := []struct {
		name string
		args args
		err  error
		want *TestBind
	}{
		{
			name: "error not nil",
			args: args{
				req:    &http.Request{Method: http.MethodPost},
				target: &TestBind{},
			},
			err:  errors.New("missing form body"),
			want: nil,
		},
		{
			name: "error is nil",
			args: args{
				req: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded; param=value"}},
					Body:   io.NopCloser(strings.NewReader("name=cnf&url=https://go-cnf.dev/")),
				},
				target: &TestBind{},
			},
			err:  nil,
			want: &TestBind{"cnf", "https://go-cnf.dev/"},
		},
		{
			name: "error BadRequest",
			args: args{
				req: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded; param=value"}},
					Body:   io.NopCloser(strings.NewReader("age=a")),
				},
				target: &TestBind2{},
			},
			err:  cnferror.BadRequest("CODEC", "Field Namespace:age ERROR:Invalid Integer Value 'a' Type 'int' Namespace 'age'"),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindForm(tt.args.req, tt.args.target)
			if !reflect.DeepEqual(err, tt.err) {
				t.Fatalf("BindForm() error = %v, err %v", err, tt.err)
			}
			if err == nil && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindForm() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
