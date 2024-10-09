package db

import (
	"gorm.io/gorm/schema"
	"reflect"
	"testing"
	"time"
)

type TestTable struct {
	ID       uint      `gorm:"id;primaryKey;autoIncrement"`
	Name     string    `gorm:"name"`
	CreateAt time.Time `gorm:"create_at"`
	UpdateAt time.Time `gorm:"update_at"`
}

func (TestTable) TableName() string {
	return "t_test_table"
}

type TestTableMapper struct {
	*BasicMapper[TestTable]
}

func NewTestTableMapper(ds *DatabaseSource) (*TestTableMapper, error) {
	bm, err := NewBasicMapper[TestTable](ds)
	if err != nil {
		return nil, err
	}
	return &TestTableMapper{
		bm,
	}, nil
}

func TestBasicMapper_Delete(t *testing.T) {
	type args[T schema.Tabler] struct {
		t *T
	}
	type testCase[T schema.Tabler] struct {
		name    string
		d       BasicMapper[T]
		args    args[T]
		want    int64
		wantErr bool
	}
	tests := []testCase[TestTable]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Delete(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicMapper_Insert(t *testing.T) {
	type args[T schema.Tabler] struct {
		t *T
	}
	type testCase[T schema.Tabler] struct {
		name    string
		d       BasicMapper[T]
		args    args[T]
		want    int64
		wantErr bool
	}
	tests := []testCase[TestTable]{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Insert(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicMapper_Select(t *testing.T) {
	type args[T schema.Tabler] struct {
		where *T
	}
	type testCase[T schema.Tabler] struct {
		name    string
		d       BasicMapper[T]
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[TestTable]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Select(tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicMapper_Update(t *testing.T) {
	type args[T schema.Tabler] struct {
		old *T
		new *T
	}
	type testCase[T schema.Tabler] struct {
		name    string
		d       BasicMapper[T]
		args    args[T]
		want    int64
		wantErr bool
	}
	tests := []testCase[TestTable]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Update(tt.args.old, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBasicMapper(t *testing.T) {
	type args struct {
		ds *DatabaseSource
	}
	type testCase[T schema.Tabler] struct {
		name    string
		args    args
		want    *BasicMapper[T]
		wantErr bool
	}
	ds, err := NewDatabaseSource(&Config{
		"/tmp/listen-tube-unit-test.db",
	})
	if err != nil {
		t.Fatalf("NewDatabaseSource() failed, err:%v", err)
	}
	tests := []testCase[TestTable]{
		{
			name: "OK",
			args: args{
				ds: ds,
			},
			want: &BasicMapper[TestTable]{
				ds,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBasicMapper[TestTable](tt.args.ds)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBasicMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBasicMapper() got = %v, want %v", got, tt.want)
			}
		})
	}
}
