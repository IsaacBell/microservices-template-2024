package consultants_biz_test

import (
	"testing"

	consultantsV1 "microservices-template-2024/api/v1/consultants"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
)

func TestCommunicationTypeString(t *testing.T) {
	testCases := []struct {
		commType consultants_biz.CommunicationType
		expected string
	}{
		{consultants_biz.COMM_TYPE_Unknown, "unknown"},
		{consultants_biz.COMM_TYPE_FromClient, "from_client"},
		{consultants_biz.COMM_TYPE_FromAdmin, "from_admin"},
		{consultants_biz.COMM_TYPE_FromSystem, "from_system"},
		{consultants_biz.CommunicationType(""), "unknown"},
	}

	for _, tc := range testCases {
		actual := tc.commType.String()
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}

func TestCommunicationTypeToProto(t *testing.T) {
	testCases := []struct {
		commType consultants_biz.CommunicationType
		expected consultantsV1.CommunicationType
	}{
		{consultants_biz.COMM_TYPE_Unknown, consultantsV1.CommunicationType_unknown},
		{consultants_biz.COMM_TYPE_FromClient, consultantsV1.CommunicationType_from_client},
		{consultants_biz.COMM_TYPE_FromAdmin, consultantsV1.CommunicationType_from_admin},
		{consultants_biz.COMM_TYPE_FromSystem, consultantsV1.CommunicationType_from_system},
		{consultants_biz.CommunicationType("invalid"), consultantsV1.CommunicationType_unknown},
	}

	for _, tc := range testCases {
		actual := tc.commType.ToProto()
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}

func TestFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected consultants_biz.CommunicationType
		hasError bool
	}{
		{"unknown", consultants_biz.COMM_TYPE_Unknown, false},
		{"from_client", consultants_biz.COMM_TYPE_FromClient, false},
		{"from_admin", consultants_biz.COMM_TYPE_FromAdmin, false},
		{"from_system", consultants_biz.COMM_TYPE_FromSystem, false},
		{"invalid", consultants_biz.COMM_TYPE_Unknown, true},
	}

	for _, tc := range testCases {
		actual, err := consultants_biz.FromString(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !tc.hasError && err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}

func TestCommunicationTypeFromProto(t *testing.T) {
	testCases := []struct {
		input    consultantsV1.CommunicationType
		expected consultants_biz.CommunicationType
	}{
		{consultantsV1.CommunicationType_unknown, consultants_biz.COMM_TYPE_Unknown},
		{consultantsV1.CommunicationType_from_client, consultants_biz.COMM_TYPE_FromClient},
		{consultantsV1.CommunicationType_from_admin, consultants_biz.COMM_TYPE_FromAdmin},
		{consultantsV1.CommunicationType_from_system, consultants_biz.COMM_TYPE_FromSystem},
	}

	for _, tc := range testCases {
		actual := consultants_biz.CommunicationTypeFromProto(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}
