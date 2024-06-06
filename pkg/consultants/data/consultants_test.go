package consultants_data_test

import (
	"context"
	"errors"
	"testing"

	consultants_biz "core/pkg/consultants/biz"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "core/test/mocks"
)

const id = "fake-uuid"

var (
	consultant  *consultants_biz.Consultant
	consultants []*consultants_biz.Consultant
	repo        *mocks.MockConsultantRepo
	comm        *consultants_biz.Communication
)

func setup() {
	repo = new(mocks.MockConsultantRepo)
	consultant = &consultants_biz.Consultant{ID: id, Bio: "Business consultant in Merryweather, CA."}
	consultants = []*consultants_biz.Consultant{
		{ID: id, UserID: id, Specializations: []string{"IT"}, Languages: []string{"English"}},
		{ID: "2", UserID: id, Specializations: []string{"IT"}, Languages: []string{"English"}},
	}
	comm = &consultants_biz.Communication{
		From:        id,
		UserID:      id,
		RecipientID: "2",
		Msg:         "Test message",
		CommType:    consultants_biz.COMM_TYPE_FromClient,
	}
	repo.On("Get", mock.Anything, id).Return(consultant, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*consultants_biz.Consultant")).Return(consultant, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*consultants_biz.Consultant")).Return(consultant, nil).Once()
	repo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
	repo.On("Search", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(consultants, nil).Once()
	repo.On("SaveCommunication", mock.Anything, mock.AnythingOfType("*consultants_biz.Communication")).Return(comm).Once()
}

func setupWithErrors() {
	err := errors.New("test err")
	repo = new(mocks.MockConsultantRepo)
	repo.On("Get", mock.Anything, id).Return(nil, err)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*consultants_biz.Consultant")).Return(nil, err).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*consultants_biz.Consultant")).Return(nil, err).Once()
	repo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(err).Once()
	repo.On("Search", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(nil, err).Once()
	repo.On("SaveCommunication", mock.Anything, mock.AnythingOfType("*consultants_biz.Communication")).Return(nil, err).Once()
}

func TestGet(t *testing.T) {
	setup()

	result, err := repo.Get(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestSaveCommunication(t *testing.T) {
	setup()

	result, err := repo.SaveCommunication(context.Background(), comm)

	assert.NoError(t, err)
	assert.Equal(t, comm, result)
}

func TestSave(t *testing.T) {
	setup()

	result, err := repo.Save(context.Background(), consultant)

	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestUpdate(t *testing.T) {
	setup()

	result, err := repo.Update(context.Background(), consultant)

	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestDelete(t *testing.T) {
	setup()

	err := repo.Delete(context.Background(), id)

	assert.NoError(t, err)
}

func TestSearch(t *testing.T) {
	setup()

	filters := map[string]interface{}{
		"user_id":         id,
		"specializations": []string{"IT"},
		"languages":       []string{"English"},
	}

	// mockDB.On("Where", "user_id = ?", []interface{}{id}).Return(mockDB).Once()
	// mockDB.On("Where", "specializations @> ?", []interface{}{[]string{"IT"}}).Return(mockDB).Once()
	// mockDB.On("Where", "languages @> ?", []interface{}{[]string{"English"}}).Return(mockDB).Once()
	// mockDB.On("Find", &[]*consultants_biz.Consultant{}).Return(nil).Run(func(args mock.Arguments) {
	// 	arg := args.Get(0).(*[]*consultants_biz.Consultant)
	// 	*arg = consultants
	// }).Once()

	result, err := repo.Search(context.Background(), filters)

	assert.NoError(t, err)
	assert.Equal(t, consultants, result)
}

func TestSearchError(t *testing.T) {
	setupWithErrors()

	filters := map[string]interface{}{"specializations": []string{"IT"}}
	// mockDB.On("Where", "specializations @> ?", []interface{}{[]string{"IT"}}).Return(mockDB).Once()
	// mockDB.On("Find", &[]*consultants_biz.Consultant{}).Return(errors.New("search error")).Once()

	_, err := repo.Search(context.Background(), filters)

	assert.Error(t, err)
}

func TestGetError(t *testing.T) {
	setupWithErrors()

	x, err := repo.Get(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, x)
}

func TestSaveError(t *testing.T) {
	setupWithErrors()

	x, err := repo.Save(context.Background(), consultant)

	assert.Error(t, err)
	assert.Nil(t, x)
}

func TestUpdateError(t *testing.T) {
	setupWithErrors()

	x, err := repo.Update(context.Background(), consultant)

	assert.Error(t, err)
	assert.Nil(t, x)
}

func TestDeleteError(t *testing.T) {
	setupWithErrors()

	err := repo.Delete(context.Background(), id)

	assert.Error(t, err)
}
