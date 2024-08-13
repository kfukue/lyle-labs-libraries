package tax

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                   //1
	"uuid",                 //2
	"name",                 //3
	"alternate_name",       //4
	"start_date",           //5
	"end_date",             //6
	"start_block",          //7
	"end_block",            //8
	"tax_rate",             //9
	"tax_rate_type_id",     //10
	"contract_address_str", //11
	"contract_address_id",  //12
	"tax_type_id",          //13
	"description",          //14
	"created_by",           //15
	"created_at",           //16
	"updated_by",           //17
	"updated_at",           //18
}
var DBColumnsInsertTaxes = []string{
	"uuid",                 //1
	"name",                 //2
	"alternate_name",       //3
	"start_date",           //4
	"end_date",             //5
	"start_block",          //6
	"end_block",            //7
	"tax_rate",             //8
	"tax_rate_type_id",     //9
	"contract_address_str", //10
	"contract_address_id",  //11
	"tax_type_id",          //12
	"description",          //13
	"created_by",           //14
	"created_at",           //15
	"updated_by",           //16
	"updated_at",           //17
}

var TestData1 = Tax{
	ID:                 utils.Ptr[int](1),                                     //1
	UUID:               "01ef85e8-2c26-441e-8c7f-71d79518ad72",                //2
	Name:               "HAMS Contract Fee",                                   //3
	AlternateName:      "HAMS Contract Fee",                                   //4
	StartDate:          utils.SampleCreatedAtTime,                             //5
	EndDate:            utils.SampleCreatedAtTime,                             //6
	StartBlock:         utils.Ptr[int](17662040),                              //7
	EndBlock:           nil,                                                   //8
	TaxRate:            utils.Ptr[decimal.Decimal](decimal.NewFromFloat(5.0)), //9
	TaxRateTypeID:      utils.Ptr[int](2),                                     //10
	ContractAddressStr: "0x48c87cDacb6Bb6BF6E5Cd85D8ee5C847084c7410",          //11
	ContractAddressID:  utils.Ptr[int](2),                                     //12
	TaxTypeID:          utils.Ptr[int](1),                                     //13
	Description:        "Buy and Sell Have 5 percent fee",                     //14
	CreatedBy:          "SYSTEM",                                              //15
	CreatedAt:          utils.SampleCreatedAtTime,                             //16
	UpdatedBy:          "SYSTEM",                                              //17
	UpdatedAt:          utils.SampleCreatedAtTime,                             //18

}

var TestData2 = Tax{
	ID:                 utils.Ptr[int](2),                                      //1
	UUID:               "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                 //2
	Name:               "JP Sales Tax",                                         //3
	AlternateName:      "Japan Sales Tax",                                      //4
	StartDate:          utils.SampleCreatedAtTime,                              //5
	EndDate:            utils.SampleCreatedAtTime,                              //6
	StartBlock:         utils.Ptr[int](1),                                      //7
	EndBlock:           nil,                                                    //8
	TaxRate:            utils.Ptr[decimal.Decimal](decimal.NewFromFloat(10.0)), //9
	TaxRateTypeID:      utils.Ptr[int](2),                                      //10
	ContractAddressStr: "",                                                     //11
	ContractAddressID:  nil,                                                    //12
	TaxTypeID:          utils.Ptr[int](3),                                      //13
	Description:        "",                                                     //14
	CreatedBy:          "SYSTEM",                                               //15
	CreatedAt:          utils.SampleCreatedAtTime,                              //16
	UpdatedBy:          "SYSTEM",                                               //17
	UpdatedAt:          utils.SampleCreatedAtTime,                              //18
}
var TestAllData = []Tax{TestData1, TestData2}

func AddTaxToMockRows(mock pgxmock.PgxPoolIface, dataList []Tax) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                 //1
			data.UUID,               //2
			data.Name,               //3
			data.AlternateName,      //4
			data.StartDate,          //5
			data.EndDate,            //6
			data.StartBlock,         //7
			data.EndBlock,           //8
			data.TaxRate,            //9
			data.TaxRateTypeID,      //10
			data.ContractAddressStr, //11
			data.ContractAddressID,  //12
			data.TaxTypeID,          //13
			data.Description,        //14
			data.CreatedBy,          //15
			data.CreatedAt,          //16
			data.UpdatedBy,          //17
			data.UpdatedAt,          //18
		)
	}
	return rows
}

func TestGetTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []Tax{targetData}
	taxID := targetData.ID
	mockRows := AddTaxToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(*taxID).WillReturnRows(mockRows)
	foundTax, err := GetTax(mock, taxID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTax", err)
	}
	if cmp.Equal(*foundTax, targetData) == false {
		t.Errorf("Expected Tax From Method GetTax: %v is different from actual %v", foundTax, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(taxID).WillReturnRows(noRows)
	foundTax, err := GetTax(mock, &taxID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTax", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTax: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxForRowErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(taxID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTax, err := GetTax(mock, &taxID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTax", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTax: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(taxID).WillReturnRows(differentModelRows)
	foundTax, err := GetTax(mock, &taxID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTax", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTax: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	taxID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM taxes").WithArgs(*taxID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveTax(mock, taxID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveTax(mock, &taxID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM taxes").WithArgs(taxID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveTax(mock, &taxID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTaxToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnRows(mockRows)
	ids := []int{1}
	foundTaxList, err := GetTaxes(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetTaxes", err)
	}
	for i, sourceTax := range dataList {
		if cmp.Equal(sourceTax, foundTaxList[i]) == false {
			t.Errorf("Expected Tax From Method GetTaxes: %v is different from actual %v", sourceTax, foundTaxList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	ids := []int{1}
	foundTaxList, err := GetTaxes(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxes", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTaxes: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnRows(differentModelRows)
	ids := []int{1}
	foundTax, err := GetTaxes(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxes", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTaxes: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTaxToMockRows(mock, dataList)
	assetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(assetID).WillReturnRows(mockRows)
	foundTaxList, err := GetTaxesByAssetID(mock, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetTaxesByAssetID", err)
	}
	for i, sourceTax := range dataList {
		if cmp.Equal(sourceTax, foundTaxList[i]) == false {
			t.Errorf("Expected Tax From Method GetTaxesByAssetID: %v is different from actual %v", sourceTax, foundTaxList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByAssetIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTaxList, err := GetTaxesByAssetID(mock, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxesByAssetID", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTaxesByAssetID: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByAssetIDForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := 1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(assetID).WillReturnRows(differentModelRows)
	foundTax, err := GetTaxesByAssetID(mock, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxesByAssetID", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTaxesByAssetID: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddTaxToMockRows(mock, dataList)
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(pq.Array(UUIDList)).WillReturnRows(mockRows)
	foundTaxList, err := GetTaxesByUUIDs(mock, UUIDList)
	if err != nil {
		t.Fatalf("an error '%s' in GetTaxesByUUIDs", err)
	}
	for i, sourceTax := range dataList {
		if cmp.Equal(sourceTax, foundTaxList[i]) == false {
			t.Errorf("Expected Tax From Method GetTaxesByUUIDs: %v is different from actual %v", sourceTax, foundTaxList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(pq.Array(UUIDList)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTaxList, err := GetTaxesByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxesByUUIDs", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTaxesByUUIDs: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxesByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	UUIDList := []string{TestData1.UUID, TestData2.UUID}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WithArgs(pq.Array(UUIDList)).WillReturnRows(differentModelRows)
	foundTax, err := GetTaxesByUUIDs(mock, UUIDList)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxesByUUIDs", err)
	}
	if foundTax != nil {
		t.Errorf("Expected Tax From Method GetTaxesByUUIDs: to be empty but got this: %v", foundTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE taxes").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.StartDate,          //3
		targetData.EndDate,            //4
		targetData.StartBlock,         //5
		targetData.EndBlock,           //6
		targetData.TaxRate,            //7
		targetData.TaxRateTypeID,      //8
		targetData.ContractAddressStr, //9
		targetData.ContractAddressID,  //10
		targetData.TaxTypeID,          //11
		targetData.Description,        //12
		targetData.UpdatedBy,          //13
		targetData.ID,                 //14
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateTax(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTaxOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE taxes").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.StartDate,          //3
		targetData.EndDate,            //4
		targetData.StartBlock,         //5
		targetData.EndBlock,           //6
		targetData.TaxRate,            //7
		targetData.TaxRateTypeID,      //8
		targetData.ContractAddressStr, //9
		targetData.ContractAddressID,  //10
		targetData.TaxTypeID,          //11
		targetData.Description,        //12
		targetData.UpdatedBy,          //13
		targetData.ID,                 //14
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO taxes").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.StartDate,          //3
		targetData.EndDate,            //4
		targetData.StartBlock,         //5
		targetData.EndBlock,           //6
		targetData.TaxRate,            //7
		targetData.TaxRateTypeID,      //8
		targetData.ContractAddressStr, //9
		targetData.ContractAddressID,  //10
		targetData.TaxTypeID,          //11
		targetData.Description,        //12
		targetData.CreatedBy,          //13
	).WillReturnRows(pgxmock.NewRows([]string{"tax_id"}).AddRow(1))
	mock.ExpectCommit()
	taxID, err := InsertTax(mock, &targetData)
	if taxID < 0 {
		t.Fatalf("taxID should not be negative ID: %d", taxID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO taxes").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.StartDate,          //3
		targetData.EndDate,            //4
		targetData.StartBlock,         //5
		targetData.EndBlock,           //6
		targetData.TaxRate,            //7
		targetData.TaxRateTypeID,      //8
		targetData.ContractAddressStr, //9
		targetData.ContractAddressID,  //10
		targetData.TaxTypeID,          //11
		targetData.Description,        //12
		targetData.CreatedBy,          //13
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	taxID, err := InsertTax(mock, &targetData)
	if taxID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error taxID: %d", taxID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTaxOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO taxes").WithArgs(
		targetData.Name,               //1
		targetData.AlternateName,      //2
		targetData.StartDate,          //3
		targetData.EndDate,            //4
		targetData.StartBlock,         //5
		targetData.EndBlock,           //6
		targetData.TaxRate,            //7
		targetData.TaxRateTypeID,      //8
		targetData.ContractAddressStr, //9
		targetData.ContractAddressID,  //10
		targetData.TaxTypeID,          //11
		targetData.Description,        //12
		targetData.CreatedBy,          //13
	).WillReturnRows(pgxmock.NewRows([]string{"tax_id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	taxID, err := InsertTax(mock, &targetData)
	if taxID >= 0 {
		t.Fatalf("Expecting -1 for taxID because of error taxID: %d", taxID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTaxes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"taxes"}, DBColumnsInsertTaxes)
	targetData := TestAllData
	err = InsertTaxes(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertTaxes", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertTaxesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"taxes"}, DBColumnsInsertTaxes).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertTaxes(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddTaxToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_rate_type_id = 1", "tax_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnRows(mockRows)
	foundTaxList, err := GetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetTaxListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundTaxList[i]) == false {
			t.Errorf("Expected sourceData From Method GetTaxListByPagination: %v is different from actual %v", sourceData, foundTaxList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundTaxList, err := GetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxListByPagination", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTaxListByPagination: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTaxListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_type_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM taxes").WillReturnRows(differentModelRows)
	foundTaxList, err := GetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTaxListByPagination", err)
	}
	if len(foundTaxList) != 0 {
		t.Errorf("Expected From Method GetTaxListByPagination: to be empty but got this: %v", foundTaxList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTaxCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM taxes").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTaxCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTaxCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTaxCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalMinersForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTaxCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTaxCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTaxCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
