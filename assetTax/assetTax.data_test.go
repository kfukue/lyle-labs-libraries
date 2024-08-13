package assettax

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var columns = []string{
	"tax_id",            //1
	"asset_id",          //2
	"uuid",              //3
	"name",              //4
	"alternate_name",    //5
	"tax_rate_override", //6
	"tax_rate_type_id",  //7
	"description",       //8
	"created_by",        //9
	"created_at",        //10
	"updated_by",        //11
	"updated_at",        //12
}

var DBColumnsInsertAssetTaxes = []string{
	"tax_id",            //1
	"asset_id",          //2
	"uuid",              //3
	"name",              //4
	"alternate_name",    //5
	"tax_rate_override", //6
	"tax_rate_type_id",  //7
	"description",       //8
	"created_by",        //9
	"created_at",        //10
	"updated_by",        //11
	"updated_at",        //12
}

var TestData1 = AssetTax{
	TaxID:           utils.Ptr[int](1),
	AssetID:         utils.Ptr[int](1),
	UUID:            "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:            "ETHER US",
	AlternateName:   "Ether United States Of America",
	TaxRateOverride: utils.Ptr[decimal.Decimal](decimal.NewFromInt(5.0)),
	TaxRateTypeID:   utils.Ptr[int](1),
	Description:     "",
	CreatedBy:       "SYSTEM",
	CreatedAt:       utils.SampleCreatedAtTime,
	UpdatedBy:       "SYSTEM",
	UpdatedAt:       utils.SampleCreatedAtTime,
}

var TestData2 = AssetTax{
	TaxID:           utils.Ptr[int](2),
	AssetID:         utils.Ptr[int](2),
	UUID:            "880607ab-2833-4ad7-a231-b983a61c7b334",
	Name:            "BTC JP",
	AlternateName:   "Bitcoin Japan",
	TaxRateOverride: utils.Ptr[decimal.Decimal](decimal.NewFromInt(10.0)),
	TaxRateTypeID:   utils.Ptr[int](2),
	Description:     "",
	CreatedBy:       "SYSTEM",
	CreatedAt:       utils.SampleCreatedAtTime,
	UpdatedBy:       "SYSTEM",
	UpdatedAt:       utils.SampleCreatedAtTime,
}
var TestAllData = []AssetTax{TestData1, TestData2}

func AddAssetTaxToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetTax) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.TaxID,           //1
			data.AssetID,         //2
			data.UUID,            //3
			data.Name,            //4
			data.AlternateName,   //5
			data.TaxRateOverride, //6
			data.TaxRateTypeID,   //8
			data.Description,     //7
			data.CreatedBy,       //9
			data.CreatedAt,       //10
			data.UpdatedBy,       //11
			data.UpdatedAt,       //12
		)
	}
	return rows
}
func TestGetAllAssetTaxesByTaxType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxTypeID := 1
	dataList := []AssetTax{TestData1}
	mockRows := AddAssetTaxToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxTypeID).WillReturnRows(mockRows)
	foundAssetTaxes, err := GetAllAssetTaxesByTaxType(mock, &taxTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllAssetTaxesByTaxType", err)
	}
	testAssetTaxes := TestAllData
	for i, foundAssetTax := range foundAssetTaxes {
		if cmp.Equal(foundAssetTax, testAssetTaxes[i]) == false {
			t.Errorf("Expected AssetTax From Method GetAllAssetTaxesByTaxType: %v is different from actual %v", foundAssetTax, testAssetTaxes[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllAssetTaxesByTaxTypeForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxTypeID := 1
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetTaxes, err := GetAllAssetTaxesByTaxType(mock, &taxTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllAssetTaxesByTaxType", err)
	}
	if len(foundAssetTaxes) != 0 {
		t.Errorf("Expected From Method GetAllAssetTaxesByTaxType: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllAssetTaxesByTaxTypeForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	taxTypeID := 1
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxTypeID).WillReturnRows(differentModelRows)
	foundAssetTaxes, err := GetAllAssetTaxesByTaxType(mock, &taxTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllAssetTaxesByTaxType", err)
	}
	if foundAssetTaxes != nil {
		t.Errorf("Expected foundAssetTaxes From Method GetAllAssetTaxesByTaxType: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []AssetTax{targetData}
	taxID := targetData.TaxID
	assetID := targetData.AssetID
	mockRows := AddAssetTaxToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(*taxID, *assetID).WillReturnRows(mockRows)
	foundAssetTax, err := GetAssetTax(mock, taxID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetTax", err)
	}
	if cmp.Equal(*foundAssetTax, targetData) == false {
		t.Errorf("Expected AssetTax From Method GetAssetTax: %v is different from actual %v", foundAssetTax, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := 999
	assetID := 999
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxID, assetID).WillReturnRows(noRows)
	foundAssetTax, err := GetAssetTax(mock, &taxID, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetTax", err)
	}
	if foundAssetTax != nil {
		t.Errorf("Expected AssetTax From Method GetAssetTax: to be empty but got this: %v", foundAssetTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxID, assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetTax, err := GetAssetTax(mock, &taxID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTax", err)
	}
	if foundAssetTax != nil {
		t.Errorf("Expected AssetTax From Method GetAssetTax: to be empty but got this: %v", foundAssetTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	taxID := -1
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WithArgs(taxID, assetID).WillReturnRows(differentModelRows)
	foundAssetTax, err := GetAssetTax(mock, &taxID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTax", err)
	}
	if foundAssetTax != nil {
		t.Errorf("Expected foundAssetTax From Method GetAssetTax: to be empty but got this: %v", foundAssetTax)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	taxID := targetData.TaxID
	assetID := targetData.AssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_taxes").WithArgs(*taxID, *assetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAssetTax(mock, taxID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAssetTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	assetID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveAssetTax(mock, &taxID, &assetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	taxID := -1
	assetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_taxes").WithArgs(taxID, assetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAssetTax(mock, &taxID, &assetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []AssetTax{TestData1, TestData2}
	mockRows := AddAssetTaxToMockRows(mock, dataList)
	assetIds := []int{1, 2}
	taxIds := []int{1, 2}
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnRows(mockRows)
	foundAssetTaxes, err := GetAssetTaxList(mock, assetIds, taxIds)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetTaxList", err)
	}
	testAssetTaxes := TestAllData
	for i, foundAssetTax := range foundAssetTaxes {
		if cmp.Equal(foundAssetTax, testAssetTaxes[i]) == false {
			t.Errorf("Expected AssetTax From Method GetAssetTaxList: %v is different from actual %v", foundAssetTax, testAssetTaxes[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIds := []int{-1, -2}
	taxIds := []int{1}
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetTaxes, err := GetAssetTaxList(mock, assetIds, taxIds)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTaxList", err)
	}
	if len(foundAssetTaxes) != 0 {
		t.Errorf("Expected From Method GetAssetTaxList: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	assetIds := []int{-1, -2}
	taxIds := []int{1}
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnRows(differentModelRows)
	foundAssetTaxes, err := GetAssetTaxList(mock, assetIds, taxIds)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTaxList", err)
	}
	if foundAssetTaxes != nil {
		t.Errorf("Expected foundAssetTaxes From Method GetAssetTaxList: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE asset_taxes").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.TaxRateOverride, //3
		targetData.TaxRateTypeID,   //4
		targetData.Description,     //5
		targetData.UpdatedBy,       //6
		targetData.TaxID,           //7
		targetData.AssetID,         //8
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAssetTax(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAssetTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetTaxOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.TaxID = nil
	err = UpdateAssetTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAssetTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateAssetTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.AssetID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE asset_taxes").WithArgs(
		targetData.Name,            //1
		targetData.AlternateName,   //2
		targetData.TaxRateOverride, //3
		targetData.TaxRateTypeID,   //4
		targetData.Description,     //5
		targetData.UpdatedBy,       //6
		targetData.TaxID,           //7
		targetData.AssetID,         //8
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateAssetTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetTax(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_taxes").WithArgs(
		targetData.TaxID,           //1
		targetData.AssetID,         //2
		targetData.Name,            //3
		targetData.AlternateName,   //4
		targetData.TaxRateOverride, //5
		targetData.TaxRateTypeID,   //6
		targetData.Description,     //7
		targetData.CreatedBy,       //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	assetID, taxID, err := InsertAssetTax(mock, &targetData)
	if assetID < 0 {
		t.Fatalf("assetID should not be negative ID: %d", assetID)
	}
	if taxID < 0 {
		t.Fatalf("taxID should not be negative ID: %d", taxID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetTax", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetTaxOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertAssetTax(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetTaxOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_taxes").WithArgs(
		targetData.TaxID,           //1
		targetData.AssetID,         //2
		targetData.Name,            //3
		targetData.AlternateName,   //4
		targetData.TaxRateOverride, //5
		targetData.TaxRateTypeID,   //6
		targetData.Description,     //7
		targetData.CreatedBy,       //8
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	assetID, taxID, err := InsertAssetTax(mock, &targetData)
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", assetID)
	}
	if taxID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", taxID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetTaxOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_taxes").WithArgs(
		targetData.TaxID,           //1
		targetData.AssetID,         //2
		targetData.Name,            //3
		targetData.AlternateName,   //4
		targetData.TaxRateOverride, //5
		targetData.TaxRateTypeID,   //6
		targetData.Description,     //7
		targetData.CreatedBy,       //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	assetID, taxID, err := InsertAssetTax(mock, &targetData)
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for assetID because of error assetID: %d", assetID)
	}
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

func TestInsertAssetTaxes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"asset_taxes"}, DBColumnsInsertAssetTaxes)
	targetData := TestAllData
	err = InsertAssetTaxes(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetTaxes", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetTaxesOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"asset_taxes"}, DBColumnsInsertAssetTaxes).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertAssetTaxes(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddAssetTaxToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_id = 1", "asset_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnRows(mockRows)
	foundAssetTaxes, err := GetAssetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetTaxListByPagination", err)
	}
	testAssetTaxes := dataList
	for i, foundAssetTax := range foundAssetTaxes {
		if cmp.Equal(foundAssetTax, testAssetTaxes[i]) == false {
			t.Errorf("Expected AssetTax From Method GetAssetTaxListByPagination: %v is different from actual %v", foundAssetTax, testAssetTaxes[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetTaxes, err := GetAssetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTaxListByPagination", err)
	}
	if len(foundAssetTaxes) != 0 {
		t.Errorf("Expected From Method GetAssetTaxListByPagination: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetTaxListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"tax_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM asset_taxes").WillReturnRows(differentModelRows)
	foundAssetTaxes, err := GetAssetTaxListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetTaxListByPagination", err)
	}
	if foundAssetTaxes != nil {
		t.Errorf("Expected From Method GetAssetTaxListByPagination: to be empty but got this: %v", foundAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetTaxCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfAssetTaxesExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM asset_taxes").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfAssetTaxesExpected))
	numOfAssetTaxes, err := GetTotalAssetTaxCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalAssetTaxCount", err)
	}
	if *numOfAssetTaxes != numOfAssetTaxesExpected {
		t.Errorf("Expected AssetTax From Method GetTotalAssetTaxCount: %d is different from actual %d", numOfAssetTaxesExpected, *numOfAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetTaxCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM asset_taxes").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfAssetTaxes, err := GetTotalAssetTaxCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalAssetTaxCount", err)
	}
	if numOfAssetTaxes != nil {
		t.Errorf("Expected numOfAssetTaxes From Method GetTotalAssetTaxCount to be empty but got this: %v", numOfAssetTaxes)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
