package tracing

import (
	"context"
	"reflect"
	"testing"

	"github.com/chnxq/xkitmod/metadata"

	"go.opentelemetry.io/otel/propagation"
)

func TestMetadata_Inject(t *testing.T) {
	type args struct {
		appName string
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "https://xkit.dev",
			args: args{"https://xkit.dev", propagation.HeaderCarrier{}},
			want: "https://xkit.dev",
		},
		{
			name: "https://github.com/xkit/xkit",
			args: args{"https://github.com/xkit/xkit", propagation.HeaderCarrier{"mode": []string{"test"}}},
			want: "https://github.com/xkit/xkit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//a := XGoKit.New(XGoKit.Name(tt.args.appName))		//TODO: temporary solution 112
			//ctx := XGoKit.NewContext(context.Background(), a)
			//m := new(Metadata)
			//m.Inject(ctx, tt.args.carrier)
			//if res := tt.args.carrier.Get(serviceHeader); tt.want != res {
			//	t.Errorf("Get(serviceHeader) :%s want: %s", res, tt.want)
			//}
		})
	}
}

func TestMetadata_Extract(t *testing.T) {
	type args struct {
		parent  context.Context
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name  string
		args  args
		want  string
		crash bool
	}{
		{
			name: "https://xkit.dev",
			args: args{
				parent:  context.Background(),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://xkit.dev"}},
			},
			want: "https://xkit.dev",
		},
		{
			name: "https://github.com/xkit/xkit",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.Metadata{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://github.com/xkit/xkit"}},
			},
			want: "https://github.com/xkit/xkit",
		},
		{
			name: "https://github.com/xkit/xkit",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.Metadata{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": nil},
			},
			crash: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Metadata{}
			ctx := b.Extract(tt.args.parent, tt.args.carrier)
			md, ok := metadata.FromServerContext(ctx)
			if !ok {
				if tt.crash {
					return
				}
				t.Errorf("expect %v, got %v", true, ok)
			}
			if !reflect.DeepEqual(md.Get(serviceHeader), tt.want) {
				t.Errorf("expect %v, got %v", tt.want, md.Get(serviceHeader))
			}
		})
	}
}

func TestFields(t *testing.T) {
	b := Metadata{}
	if !reflect.DeepEqual(b.Fields(), []string{"x-md-service-name"}) {
		t.Errorf("expect %v, got %v", []string{"x-md-service-name"}, b.Fields())
	}
}
