package consultants_service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	communicationsV1 "microservices-template-2024/api/v1/communications"
	consultantsV1 "microservices-template-2024/api/v1/consultants"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	consultants_service "microservices-template-2024/pkg/consultants/service"
	mocks "microservices-template-2024/test/mocks"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const id = "fake-uuid"
const yoe int32 = 5

var (
	logger     log.Logger
	repo       *mocks.MockConsultantRepo
	action     *consultants_biz.ConsultantAction
	consultant *consultants_biz.Consultant
	service    *consultants_service.ConsultantService

	comm *consultants_biz.Communication
)

func setup() {
	repo = new(mocks.MockConsultantRepo)
	logger = log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
	action = consultants_biz.NewConsultantAction(repo, logger)
	consultant = &consultants_biz.Consultant{ID: id, YearsOfExperience: 5}
	comm = &consultants_biz.Communication{
		From: id,
		CommType: consultants_biz.CommunicationType(
			communicationsV1.CommunicationType_from_client.String()),
		RecipientID: "2",
		Msg:         "Test msg",
		UserID:      id,
	}
	repo.On("Get", context.Background(), id).Return(consultant, nil)
	repo.On("Save", context.Background(), consultant).Return(consultant, nil)
	repo.On("Update", mock.Anything, consultant).Return(consultant, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)
	repo.On("SaveCommunication", mock.Anything, mock.Anything).Return(comm, nil)

	consultantList := []*consultants_biz.Consultant{{ID: "1", YearsOfExperience: 5}, {ID: "2", YearsOfExperience: 22}}
	repo.On("Search", mock.Anything, mock.Anything).Return(consultantList, nil)

	service = consultants_service.NewConsultantService(action)
}

func setupServiceWithErrors() {
	err := errors.New("Test Err")
	repo = new(mocks.MockConsultantRepo)
	logger = log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
	action = consultants_biz.NewConsultantAction(repo, logger)
	consultant = &consultants_biz.Consultant{ID: id, YearsOfExperience: 5}
	repo.On("Get", context.Background(), id).Return(nil, err)
	repo.On("Save", context.Background(), consultant).Return(nil, err)
	repo.On("Update", mock.Anything, consultant).Return(nil, err)
	repo.On("Delete", mock.Anything, id).Return(err)
	repo.On("Search", mock.Anything, mock.Anything).Return(nil, err)
	repo.On("SaveCommunication", mock.Anything, mock.Anything).Return(nil, err)

	service = consultants_service.NewConsultantService(action)
}

func TestGetConsultant(t *testing.T) {
	setup()

	req := &consultantsV1.GetConsultantRequest{Id: id}
	reply, err := service.GetConsultant(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, id, reply.Consultant.Id)
	assert.Equal(t, yoe, reply.Consultant.YearsOfExperience)
}

func TestCreateConsultant(t *testing.T) {
	setup()

	req := &consultantsV1.CreateConsultantRequest{Consultant: consultants_biz.ConsultantToProtoData(consultant)}
	reply, err := service.CreateConsultant(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, id, reply.Consultant.Id)
	assert.Equal(t, yoe, reply.Consultant.YearsOfExperience)
}

func TestCreateConsultantError(t *testing.T) {
	setupServiceWithErrors()

	req := &consultantsV1.CreateConsultantRequest{Consultant: consultants_biz.ConsultantToProtoData(consultant)}
	reply, err := service.CreateConsultant(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestUpdateConsultant(t *testing.T) {
	setup()

	req := &consultantsV1.UpdateConsultantRequest{Consultant: consultants_biz.ConsultantToProtoData(consultant)}
	reply, err := service.UpdateConsultant(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, id, reply.Consultant.Id)
	assert.Equal(t, yoe, reply.Consultant.YearsOfExperience)
}

func TestUpdateConsultantError(t *testing.T) {
	setupServiceWithErrors()

	req := &consultantsV1.UpdateConsultantRequest{Consultant: consultants_biz.ConsultantToProtoData(consultant)}
	reply, err := service.UpdateConsultant(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteConsultant(t *testing.T) {
	setup()

	req := &consultantsV1.DeleteConsultantRequest{Id: id}
	reply, err := service.DeleteConsultant(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
}

func TestDeleteConsultantError(t *testing.T) {
	setupServiceWithErrors()

	req := &consultantsV1.DeleteConsultantRequest{Id: id}
	reply, err := service.DeleteConsultant(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestListConsultants(t *testing.T) {
	setup()

	req := &consultantsV1.ListConsultantsRequest{Specializations: []string{"IT"}}
	reply, err := service.ListConsultants(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Len(t, reply.Consultants, 2)
	assert.Equal(t, "1", reply.Consultants[0].Id)
	assert.Equal(t, yoe, reply.Consultants[0].YearsOfExperience)
	assert.Equal(t, "2", reply.Consultants[1].Id)
	assert.Equal(t, int32(22), reply.Consultants[1].YearsOfExperience)
}

func TestListConsultantsError(t *testing.T) {
	setupServiceWithErrors()

	req := &consultantsV1.ListConsultantsRequest{Specializations: []string{"IT"}}
	reply, err := service.ListConsultants(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestSendComm(t *testing.T) {
	setup()

	c := consultants_biz.CommunicationToProtoData(comm)
	req := &consultantsV1.SendCommsRequest{Comm: c}
	reply, err := service.SendComm(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, id, reply.Ack.UserId)
	assert.Equal(t, "Test msg", reply.Ack.Msg)
	assert.Equal(t, int32(1), reply.Ack.RecipientsCount)
	assert.Equal(t, []string{"2"}, reply.Ack.Recipients)
}

func TestSendCommError(t *testing.T) {
	setupServiceWithErrors()

	req := &consultantsV1.SendCommsRequest{Comm: consultants_biz.CommunicationToProtoData(comm)}
	reply, err := service.SendComm(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, reply)
}
