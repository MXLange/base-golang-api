package errors

import (
	"reflect"
	"testing"
)

func Test_New(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want *AppError
	}{
		{
			name: "should create a new AppError with the given code",
			args: args{code: 404},
			want: &AppError{code: 404},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewError(t *testing.T) {
	type args struct {
		field   string
		message string
	}
	tests := []struct {
		name string
		args args
		want *FieldError
	}{
		{
			name: "should create a new error with the given field and message",
			args: args{
				field:   "name",
				message: "name is required",
			},
			want: &FieldError{
				field:   "name",
				message: "name is required",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.field, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appError_MarshalJSON(t *testing.T) {
	type fields struct {
		code   int
		errors []FieldError
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "should marshal the AppError to JSON",
			fields: fields{
				code: 400,
				errors: []FieldError{
					{field: "name", message: "name is required"},
				},
			},
			want:    []byte(`{"code":400,"errors":[{"field":"name","message":"name is required"}]}`),
			wantErr: false,
		},
		{
			name: "should marshal the AppError to JSON",
			fields: fields{
				code: 400,
				errors: []FieldError{
					{message: "name is required"},
				},
			},
			want:    []byte(`{"code":400,"errors":[{"message":"name is required"}]}`),
			wantErr: false,
		},
		{
			name: "should marshal the AppError to JSON",
			fields: fields{
				code: 400,
				errors: []FieldError{
					{field: "name"},
				},
			},
			want:    []byte(`{"code":400,"errors":[{"field":"name"}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := &AppError{
				code:   tt.fields.code,
				errors: tt.fields.errors,
			}
			got, err := ae.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppError.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_appError_WithErrors(t *testing.T) {
	type args struct {
		errors []FieldError
	}
	tests := []struct {
		name string
		ae   *AppError
		args args
		want *AppError
	}{
		{
			name: "should set the errors field of the AppError with a list of err",
			ae:   &AppError{code: 400},
			args: args{
				errors: []FieldError{
					{field: "name", message: "name is required"},
				},
			},
			want: &AppError{
				code: 400,
				errors: []FieldError{
					{field: "name", message: "name is required"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ae.WithErrors(tt.args.errors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppError.WithErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_err_WithField(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		e    *FieldError
		args args
		want *FieldError
	}{
		{
			name: "should set the field of the FieldError instance",
			e:    &FieldError{message: "name is required"},
			args: args{field: "name"},
			want: &FieldError{field: "name", message: "name is required"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.WithField(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("err.WithField() = %v, want %v", got, tt.want)
			}
		})
	}
}
