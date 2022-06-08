package visitor

import (
	"errors"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/parser"
	"reflect"
	"testing"
	"text/template"
)

func TestFindArgs_Error(t *testing.T) {
	type fields struct {
		conf   config.Config
		valMap map[string]struct{}
		err    error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "no error", fields: fields{}},
		{"error", fields{err: errors.New("error")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := FindArgs{
				conf:   tt.fields.conf,
				valMap: tt.fields.valMap,
				err:    tt.fields.err,
			}
			if err := v.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindArgs_Values(t *testing.T) {
	type fields struct {
		conf   config.Config
		valMap map[string]struct{}
		err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"simple", fields{valMap: map[string]struct{}{"a": {}}}, []string{"a="}},
		{"nested", fields{valMap: map[string]struct{}{"a.b": {}}}, []string{"a.b="}},
		{"duplicate", fields{conf: config.Config{Values: map[string]any{"b": "b"}}, valMap: map[string]struct{}{"b": {}}}, []string{}},
		{"reserved", fields{valMap: map[string]struct{}{"Value": {}}}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := FindArgs{
				conf:   tt.fields.conf,
				valMap: tt.fields.valMap,
				err:    tt.fields.err,
			}
			if got := v.Values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindArgs_Visit(t *testing.T) {
	defaultConf := config.New()
	defaultConf.Values["b"] = "b"

	tests := []struct {
		name    string
		v       map[string]struct{}
		source  string
		want    map[string]struct{}
		wantErr bool
	}{
		{"simple", make(map[string]struct{}), "a #yampl {{ .a }}", map[string]struct{}{"a": {}}, false},
		{"invalid template", make(map[string]struct{}), "a #yampl {{", map[string]struct{}{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := parser.ParseBytes([]byte(tt.source))

			v := NewFindArgs(defaultConf)
			if v.Visit(file.Docs[0].Body); (v.err != nil) != tt.wantErr {
				t.Errorf("Visitor() error = %v, wantErr %v", v.err, tt.wantErr)
				return
			}

			got := v.valMap
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Visitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFindArgs(t *testing.T) {
	type args struct {
		conf config.Config
	}
	tests := []struct {
		name string
		args args
		want FindArgs
	}{
		{"default", args{conf: config.Config{}}, FindArgs{conf: config.Config{}, valMap: make(map[string]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFindArgs(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFindArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listNodeFields(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"simple", args{"{{ .a }}"}, []string{"{{.a}}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, _ := template.New("").Parse(tt.args.source)

			if got := listNodeFields(tmpl.Tree.Root, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listNodeFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listTemplFields(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"simple", args{"{{ .a }}"}, []string{"{{.a}}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, _ := template.New("").Parse(tt.args.source)

			if got := listTemplFields(tmpl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listTemplFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
