package asset

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateLookupByContractAddressFromAssetList(t *testing.T) {
	targetData := TestAllData
	result := CreateLookupByContractAddressFromAssetList(targetData)
	if cmp.Equal(*result[strings.ToLower(TestData1.ContractAddress)], TestData1) == false {
		t.Errorf("Expected Asset From Method CreateLookupByContractAddressFromAssetList: %v is different from actual %v", result[strings.ToLower(TestData1.ContractAddress)], TestData1)
	}
}
