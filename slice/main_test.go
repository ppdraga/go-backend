package slice

import (
	"reflect"
	"testing"
)

func TestPageOfSliceStr(t *testing.T) {
	testSliceLen10 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	type args struct {
		s       []string
		page    uint32
		perPage uint32
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test 1. len=10, page=1, perPage=5",
			args: args{
				s:       testSliceLen10,
				page:    1,
				perPage: 5,
			},
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "Test 2. len=10, page=2, perPage=5",
			args: args{
				s:       testSliceLen10,
				page:    2,
				perPage: 5,
			},
			want: []string{"6", "7", "8", "9", "10"},
		},
		{
			name: "Test 3. len=10, page=3, perPage=5",
			args: args{
				s:       testSliceLen10,
				page:    3,
				perPage: 5,
			},
			want: []string{},
		},
		{
			name: "Test 4. len=10, page=1, perPage=15",
			args: args{
				s:       testSliceLen10,
				page:    1,
				perPage: 15,
			},
			want: testSliceLen10,
		},
		{
			name: "Test 5. len=10, page=5, perPage=1",
			args: args{
				s:       testSliceLen10,
				page:    5,
				perPage: 1,
			},
			want: []string{"5"},
		},
		{
			name: "Test 6. len=10, page=4, perPage=3",
			args: args{
				s:       testSliceLen10,
				page:    4,
				perPage: 3,
			},
			want: []string{"10"},
		},
		{
			name: "Test 7. len=10, page=1, perPage=0",
			args: args{
				s:       testSliceLen10,
				page:    1,
				perPage: 0,
			},
			want: []string{},
		},
		{
			name: "Test 8. len=10, page=0, perPage=5",
			args: args{
				s:       testSliceLen10,
				page:    0,
				perPage: 5,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PageOfSliceStr(tt.args.s, tt.args.page, tt.args.perPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageOfSliceStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitSliceStr(t *testing.T) {
	type args struct {
		s         []string
		batchSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Тест 14/5",
			args: args{
				s:         []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14"},
				batchSize: 5,
			},
			want: [][]string{
				[]string{"1", "2", "3", "4", "5"},
				[]string{"6", "7", "8", "9", "10"},
				[]string{"11", "12", "13", "14"},
			},
		},
		{
			name: "Тест 3/10",
			args: args{
				s:         []string{"1", "2", "3"},
				batchSize: 10,
			},
			want: [][]string{
				[]string{"1", "2", "3"},
			},
		},
		{
			name: "Тест 1/1",
			args: args{
				s:         []string{"1"},
				batchSize: 1,
			},
			want: [][]string{
				[]string{"1"},
			},
		},
		{
			name: "Тест 1/0",
			args: args{
				s:         []string{"1"},
				batchSize: 0,
			},
			want: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitSliceStr(tt.args.s, tt.args.batchSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSliceStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
