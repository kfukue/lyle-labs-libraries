package assetsource

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var columns = []string{
	"source_id",         //1
	"asset_id",          //2
	"uuid",              //3
	"name",              //4
	"alternate_name",    //5
	"source_identifier", //6
	"description",       //7
	"source_data",       //8
	"created_by",        //9
	"created_at",        //10
	"updated_by",        //11
	"updated_at",        //12
}

var data1 = AssetSource{
	SourceID:         utils.Ptr[int](1),
	AssetID:          utils.Ptr[int](1),
	UUID:             "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:             "ETHER",
	AlternateName:    "Ether",
	SourceIdentifier: "ETH",
	Description:      "",
	SourceData:       nil,
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
}

var data2 = AssetSource{
	SourceID:         utils.Ptr[int](2),
	AssetID:          utils.Ptr[int](2),
	UUID:             "880607ab-2833-4ad7-a231-b983a61c7b334",
	Name:             "BTC",
	AlternateName:    "Bitcoin",
	SourceIdentifier: "BTC",
	Description:      "",
	SourceData:       nil,
	CreatedBy:        "SYSTEM",
	CreatedAt:        utils.SampleCreatedAtTime,
	UpdatedBy:        "SYSTEM",
	UpdatedAt:        utils.SampleCreatedAtTime,
}
var allData = []AssetSource{data1, data2}

func AddAssetSourceToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetSource) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.SourceID,         //1
			data.AssetID,          //2
			data.UUID,             //3
			data.Name,             //4
			data.AlternateName,    //5
			data.SourceIdentifier, //6
			data.Description,      //7
			data.SourceData,       //8
			data.CreatedBy,        //9
			data.CreatedAt,        //10
			data.UpdatedBy,        //11
			data.UpdatedAt,        //12
		)
	}
	return rows
}
func TestGetAllAssetSourceBySourceAndAssetType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	sourceID := targetData.SourceID
	assetTypeID := 1
	dataList := allData
	mockRows := AddAssetSourceToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(*sourceID, assetTypeID).WillReturnRows(mockRows)
	foundAssetSources, err := GetAllAssetSourceBySourceAndAssetType(mock, sourceID, &assetTypeID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllAssetSourceBySourceAndAssetType", err)
	}
	testAssetSources := allData
	for i, foundAssetSource := range foundAssetSources {
		if cmp.Equal(foundAssetSource, testAssetSources[i]) == false {
			t.Errorf("Expected AssetSource From Method GetAllAssetSourceBySourceAndAssetType: %v is different from actual %v", foundAssetSource, testAssetSources[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllAssetSourceBySourceAndAssetTypeForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	sourceID := targetData.SourceID
	assetTypeID := 1
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(*sourceID, assetTypeID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetSources, err := GetAllAssetSourceBySourceAndAssetType(mock, sourceID, &assetTypeID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllAssetSourceBySourceAndAssetType", err)
	}
	if len(foundAssetSources) != 0 {
		t.Errorf("Expected From Method GetAllAssetSourceBySourceAndAssetType: to be empty but got this: %v", foundAssetSources)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []AssetSource{targetData}
	sourceID := targetData.SourceID
	assetID := targetData.AssetID
	mockRows := AddAssetSourceToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(*sourceID, *assetID).WillReturnRows(mockRows)
	foundAssetSource, err := GetAssetSource(mock, sourceID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSource", err)
	}
	if cmp.Equal(*foundAssetSource, targetData) == false {
		t.Errorf("Expected AssetSource From Method GetAssetSource: %v is different from actual %v", foundAssetSource, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	assetID := 999
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(sourceID, assetID).WillReturnRows(noRows)
	foundAssetSource, err := GetAssetSource(mock, &sourceID, &assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSource", err)
	}
	if foundAssetSource != nil {
		t.Errorf("Expected AssetSource From Method GetAssetSource: to be empty but got this: %v", foundAssetSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	assetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(sourceID, assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetSource, err := GetAssetSource(mock, &sourceID, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetSource", err)
	}
	if foundAssetSource != nil {
		t.Errorf("Expected AssetSource From Method GetAssetSource: to be empty but got this: %v", foundAssetSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceByTicker(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data2
	dataList := []AssetSource{targetData}
	sourceID := targetData.SourceID
	sourceIdentifier := targetData.SourceIdentifier
	mockRows := AddAssetSourceToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(*sourceID, sourceIdentifier).WillReturnRows(mockRows)
	foundAssetSource, err := GetAssetSourceByTicker(mock, sourceID, sourceIdentifier)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSourceByTicker", err)
	}
	if cmp.Equal(*foundAssetSource, targetData) == false {
		t.Errorf("Expected AssetSource From Method GetAssetSourceByTicker: %v is different from actual %v", foundAssetSource, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceByTickerForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	sourceIdentifier := "nothing"
	noRows := pgxmock.NewRows(columns)
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(sourceID, sourceIdentifier).WillReturnRows(noRows)
	foundAssetSource, err := GetAssetSourceByTicker(mock, &sourceID, sourceIdentifier)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSourceByTicker", err)
	}
	if foundAssetSource != nil {
		t.Errorf("Expected AssetSource From Method GetAssetSourceByTicker: to be empty but got this: %v", foundAssetSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceByTickerForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := 999
	sourceIdentifier := "nothing"
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WithArgs(sourceID, sourceIdentifier).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetSource, err := GetAssetSourceByTicker(mock, &sourceID, sourceIdentifier)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetSourceByTicker", err)
	}
	if foundAssetSource != nil {
		t.Errorf("Expected AssetSource From Method GetAssetSourceByTicker: to be empty but got this: %v", foundAssetSource)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	sourceID := targetData.SourceID
	assetID := targetData.AssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_sources").WithArgs(*sourceID, *assetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAssetSource(mock, sourceID, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAssetSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	sourceID := -1
	assetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_sources").WithArgs(sourceID, assetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAssetSource(mock, &sourceID, &assetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []AssetSource{data1, data2}
	mockRows := AddAssetSourceToMockRows(mock, dataList)
	assetIds := []int{1, 2}
	sourceIds := []int{1}
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WillReturnRows(mockRows)
	foundAssetSources, err := GetAssetSourceList(mock, assetIds, sourceIds)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSourceList", err)
	}
	testAssetSources := allData
	for i, foundAssetSource := range foundAssetSources {
		if cmp.Equal(foundAssetSource, testAssetSources[i]) == false {
			t.Errorf("Expected AssetSource From Method GetAssetSourceList: %v is different from actual %v", foundAssetSource, testAssetSources[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetIds := []int{-1, -2}
	sourceIds := []int{1}
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetSources, err := GetAssetSourceList(mock, assetIds, sourceIds)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetSourceList", err)
	}
	if len(foundAssetSources) != 0 {
		t.Errorf("Expected From Method GetAssetSourceList: to be empty but got this: %v", foundAssetSources)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE asset_sources").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.SourceIdentifier, //3
		targetData.Description,      //4
		targetData.SourceData,       //5
		targetData.UpdatedBy,        //6
		targetData.SourceID,         //7
		targetData.AssetID,          //8
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAssetSource(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAssetSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	targetData.AssetID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE asset_sources").WithArgs(
		targetData.Name,             //1
		targetData.AlternateName,    //2
		targetData.SourceIdentifier, //3
		targetData.Description,      //4
		targetData.SourceData,       //5
		targetData.UpdatedBy,        //6
		targetData.SourceID,         //7
		targetData.AssetID,          //8
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateAssetSource(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetSource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	targetData.Name = "New Name"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_sources").WithArgs(
		targetData.SourceID,         //1
		targetData.AssetID,          //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.SourceIdentifier, //5
		targetData.Description,      //6
		targetData.SourceData,       //7
		targetData.CreatedBy,        //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	assetID, sourceID, err := InsertAssetSource(mock, &targetData)
	if assetID < 0 {
		t.Fatalf("assetID should not be negative ID: %d", assetID)
	}
	if sourceID < 0 {
		t.Fatalf("sourceID should not be negative ID: %d", sourceID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetSource", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetSourceOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_sources").WithArgs(
		targetData.SourceID,         //1
		targetData.AssetID,          //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.SourceIdentifier, //5
		targetData.Description,      //6
		targetData.SourceData,       //7
		targetData.CreatedBy,        //8
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	assetID, sourceID, err := InsertAssetSource(mock, &targetData)
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", assetID)
	}
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error ID: %d", sourceID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetSourceOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := data1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_sources").WithArgs(
		targetData.SourceID,         //1
		targetData.AssetID,          //2
		targetData.Name,             //3
		targetData.AlternateName,    //4
		targetData.SourceIdentifier, //5
		targetData.Description,      //6
		targetData.SourceData,       //7
		targetData.CreatedBy,        //8
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	assetID, sourceID, err := InsertAssetSource(mock, &targetData)
	if assetID >= 0 {
		t.Fatalf("Expecting -1 for assetID because of error assetID: %d", assetID)
	}
	if sourceID >= 0 {
		t.Fatalf("Expecting -1 for sourceID because of error sourceID: %d", sourceID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := allData
	mockRows := AddAssetSourceToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"source_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WillReturnRows(mockRows)
	foundAssetSources, err := GetAssetSourceListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetSourceListByPagination", err)
	}
	testAssetSources := dataList
	for i, foundAssetSource := range foundAssetSources {
		if cmp.Equal(foundAssetSource, testAssetSources[i]) == false {
			t.Errorf("Expected AssetSource From Method GetAssetSourceListByPagination: %v is different from actual %v", foundAssetSource, testAssetSources[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAssetSourceListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"source_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM asset_sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAssetSources, err := GetAssetSourceListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAssetSourceListByPagination", err)
	}
	if len(foundAssetSources) != 0 {
		t.Errorf("Expected From Method GetAssetSourceListByPagination: to be empty but got this: %v", foundAssetSources)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetSourceCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfAssetSourcesExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM asset_sources").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfAssetSourcesExpected))
	numOfAssetSources, err := GetTotalAssetSourceCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalAssetSourceCount", err)
	}
	if *numOfAssetSources != numOfAssetSourcesExpected {
		t.Errorf("Expected AssetSource From Method GetTotalAssetSourceCount: %d is different from actual %d", numOfAssetSourcesExpected, *numOfAssetSources)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAssetSourceCountForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT COUNT(.*) FROM asset_sources").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfAssetSources, err := GetTotalAssetSourceCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalAssetSourceCount", err)
	}
	if numOfAssetSources != nil {
		t.Errorf("Expected numOfAssetSources From Method GetTotalAssetSourceCount to be empty but got this: %v", numOfAssetSources)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
