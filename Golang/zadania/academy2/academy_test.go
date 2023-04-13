package academy

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGradeStudent(t *testing.T) {
	zbyszek := Sophomore{
		name:       "Zbyszek",
		grades:     []int{1, 1, 1, 1},
		project:    1,
		attendance: []bool{false, true, false, false},
	}
	ania := Sophomore{
		name:       "Ania",
		grades:     []int{5, 5, 5, 5},
		project:    4,
		attendance: []bool{true, true, false, true},
	}
	bartek := Sophomore{
		name:       "Bartek",
		grades:     []int{5, 5, 5, 5},
		project:    1,
		attendance: []bool{true, true, false, true},
	}
	vasylij := Sophomore{
		name:       "Vasylij",
		grades:     []int{5, 5, 5, 5},
		project:    5,
		attendance: []bool{false, false, false, false},
	}

	repository := NewMockRepository(t)
	repository.On("Get", "Zbyszek").Return(zbyszek, nil)
	repository.On("Get", "Ania").Return(ania, nil)
	repository.On("Get", "Bartek").Return(bartek, nil)
	repository.On("Get", "Vasylij").Return(vasylij, nil)
	repository.On("Save", "Zbyszek", uint8(2)).Return(nil)
	repository.On("Save", "Ania", uint8(3)).Return(nil)
	repository.On("Save", "Bartek", uint8(2)).Return(nil)
	repository.On("Save", "Vasylij", uint8(2)).Return(nil)

	type args struct {
		name string
	}
	type result struct {
		year int
	}
	tests := []struct {
		name    string
		args    args
		expect  result
		wantErr bool
	}{
		{"Should not be promoted due to low grades", args{"Zbyszek"}, result{2}, false},
		{"Should not be promoted due to low project grade", args{"Bartek"}, result{2}, false},
		{"Should not be promoted due to low attendence", args{"Vasylij"}, result{2}, false},
		{"Should be promoted", args{"Ania"}, result{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GradeStudent(repository, tt.args.name)

			assert.Equal(t, tt.wantErr, err != nil)
			repository.AssertCalled(t, "Save", tt.args.name, uint8(tt.expect.year))
		})
	}
}

func TestGradeYear(t *testing.T) {
	zbyszek := Sophomore{
		name:       "Zbyszek",
		grades:     []int{1, 1, 1, 1},
		project:    1,
		attendance: []bool{false, true, false, false},
	}
	ania := Sophomore{
		name:       "Ania",
		grades:     []int{5, 5, 5, 5},
		project:    4,
		attendance: []bool{true, true, false, true},
	}

	repository := NewMockRepository(t)
	repository.On("Get", "Zbyszek").Return(zbyszek, nil)
	repository.On("Get", "Ania").Return(ania, nil)
	repository.On("Save", "Zbyszek", uint8(2)).Return(nil)
	repository.On("Save", "Ania", uint8(3)).Return(nil)
	repository.On("List", uint8(2)).Return([]string{zbyszek.Name(), ania.Name()}, nil)

	type args struct {
		year uint8
	}
	type result struct {
		yearCalled uint8
		count      int
	}
	tests := []struct {
		name    string
		args    args
		expect  []result
		wantErr bool
	}{
		{"Should once promote and once stay the same", args{2}, []result{{2, 1}, {3, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GradeYear(repository, tt.args.year)

			assert.Equal(t, tt.wantErr, err != nil)
			repository.AssertNumberOfCalls(t, "Save", len(tt.expect))
			for _, expected := range tt.expect {
				repository.AssertCalled(t, "Save", mock.AnythingOfType("string"), expected.yearCalled)
			}
		})
	}
}
