package db

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm/schema"
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

var fixedTime = time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

var testTable = &TestTable{
	ID:       1,
	Name:     "Test",
	CreateAt: fixedTime,
	UpdateAt: fixedTime,
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

func setupSuite(t *testing.T) func(t *testing.T) {
	// Setup code here
	return func(t *testing.T) {
		// Teardown: delete test db files
		files, err := filepath.Glob("/tmp/listen-tube-unit-test-*.db")
		if err != nil {
			t.Fatalf("Failed to list test db files: %v", err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				t.Fatalf("Failed to delete test db file %s: %v", file, err)
			}
		}
	}
}

func setupTest(t *testing.T, mapper *TestTableMapper) func(t *testing.T) {
	// Insert initial entities
	_, err := mapper.Insert(testTable)
	if err != nil {
		t.Fatalf("Failed to insert testTable: %v", err)
	}

	// Return a function to teardown the test
	return func(t *testing.T) {
		// Teardown: clean up the database
		mapper.DB.Exec("DELETE FROM " + testTable.TableName())
	}
}

func MockTestTableMapper() *TestTableMapper {
	conf := &Config{
		DSN: fmt.Sprintf("/tmp/listen-tube-unit-test-%d.db", time.Now().UnixNano()),
	}
	ds, err := NewDatabaseSource(conf)
	if err != nil {
		panic(err)
	}

	mapper, err := NewTestTableMapper(ds)
	if err != nil {
		panic(err)
	}

	return mapper
}

func TestBasicMapper_Delete(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			t *TestTable
		}
		want    int64
		wantErr bool
	}{
		{
			name:    "Delete existing record",
			args:    struct{ t *TestTable }{t: &TestTable{ID: 1}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Delete non-existing record",
			args:    struct{ t *TestTable }{t: &TestTable{ID: 999}},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := MockTestTableMapper()
			teardownTest := setupTest(t, mapper)
			defer teardownTest(t)

			got, err := mapper.Delete(tt.args.t)
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			t *TestTable
		}
		want    int64
		wantErr bool
	}{
		{
			name:    "Insert new record",
			args:    struct{ t *TestTable }{t: &TestTable{Name: "New Record", CreateAt: fixedTime, UpdateAt: fixedTime}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Insert duplicate record",
			args:    struct{ t *TestTable }{t: &TestTable{ID: 1, Name: "Duplicate Record", CreateAt: fixedTime, UpdateAt: fixedTime}},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := MockTestTableMapper()
			teardownTest := setupTest(t, mapper)
			defer teardownTest(t)

			got, err := mapper.Insert(tt.args.t)
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			where TestTable
		}
		want    []*TestTable
		wantErr bool
	}{
		{
			name:    "Select existing record",
			args:    struct{ where TestTable }{where: TestTable{ID: 1}},
			want:    []*TestTable{testTable},
			wantErr: false,
		},
		{
			name:    "Select non-existing record",
			args:    struct{ where TestTable }{where: TestTable{ID: 999}},
			want:    []*TestTable{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := MockTestTableMapper()
			teardownTest := setupTest(t, mapper)
			defer teardownTest(t)

			got, err := mapper.Select(&tt.args.where)
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name string
		args struct {
			old *TestTable
			new *TestTable
		}
		want    int64
		wantErr bool
	}{
		{
			name: "Update existing record",
			args: struct {
				old *TestTable
				new *TestTable
			}{old: &TestTable{ID: 1}, new: &TestTable{Name: "Updated", CreateAt: fixedTime, UpdateAt: fixedTime}},
			want:    1,
			wantErr: false,
		},
		{
			name: "Update non-existing record",
			args: struct {
				old *TestTable
				new *TestTable
			}{old: &TestTable{ID: 999}, new: &TestTable{Name: "Updated", CreateAt: fixedTime, UpdateAt: fixedTime}},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := MockTestTableMapper()
			teardownTest := setupTest(t, mapper)
			defer teardownTest(t)

			got, err := mapper.Update(tt.args.old, tt.args.new)
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
