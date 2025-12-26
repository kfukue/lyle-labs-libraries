package aiModels

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",              //1
	"uuid",            //2
	"name",            //3
	"alternate_name",  //4
	"url",             //5
	"ticker",          //6
	"description",     //7
	"ollama_name",     //8
	"params_size",     //9
	"quantiized_size", //10
	"base_model_id",   //11
	"created_by",      //12
	"created_at",      //13
	"updated_by",      //14
	"updated_at",      //15
}
var DBColumnsInsertAIModels = []string{
	"uuid",            //1
	"name",            //2
	"alternate_name",  //3
	"url",             //4
	"ticker",          //5
	"description",     //6
	"ollama_name",     //7
	"params_size",     //8
	"quantiized_size", //9
	"base_model_id",   //10
	"created_by",      //11
	"created_at",      //12
	"updated_by",      //13
	"updated_at",      //14
}

var TestData1 = AIModel{
	ID:            utils.Ptr[int](1),                      //1
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72", //2
	Name:          "Mistral 7B Instruct",                  //3
	AlternateName: "mistral-new:latest ",                  //4
	URL:           "https://wwww.coingecko.com",           //5
	Ticker:        "mistral-new:latest",                   //6
	Description:   "",                                     //7
	OllamaName:    "mistral-new:latest ",                  //8
	ParamsSize:    7000000000000,                          //9
	QuantizedSize: "Q4",                                   //10
	BaseModelID:   utils.Ptr[int](1),                      //11
	CreatedBy:     "SYSTEM",                               //12
	CreatedAt:     utils.SampleCreatedAtTime,              //13
	UpdatedBy:     "SYSTEM",                               //14
	UpdatedAt:     utils.SampleCreatedAtTime,              //15

}

var TestData2 = AIModel{
	ID:            utils.Ptr[int](2),                        //1
	UUID:          "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",   //2
	Name:          "Deep Seek R1 8B",                        //3
	AlternateName: "deepseek-r1:8b",                         //4
	URL:           "https://ollama.com/library/deepseek-r1", //5
	Ticker:        "deepseek-r1:8b",                         //6
	Description:   "",                                       //7
	OllamaName:    "deepseek-r1:8b",                         //8
	ParamsSize:    8000000000000,                            //9
	QuantizedSize: "Q4",                                     //10
	BaseModelID:   utils.Ptr[int](1),                        //11
	CreatedBy:     "SYSTEM",                                 //12
	CreatedAt:     utils.SampleCreatedAtTime,                //13
	UpdatedBy:     "SYSTEM",                                 //14
	UpdatedAt:     utils.SampleCreatedAtTime,                //15
}
var TestAllData = []AIModel{TestData1, TestData2}

func AddAIModelToMockRows(mock pgxmock.PgxPoolIface, dataList []AIModel) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.URL,           //5
			data.Ticker,        //6
			data.Description,   //7
			data.OllamaName,    //8
			data.ParamsSize,    //9
			data.QuantizedSize, //10
			data.BaseModelID,   //11
			data.CreatedBy,     //12
			data.CreatedAt,     //13
			data.UpdatedBy,     //14
			data.UpdatedAt,     //15
		)
	}
	return rows
}

func TestGetAIModel(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []AIModel{targetData}
	aiModelID := targetData.ID
	mockRows := AddAIModelToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WithArgs(*aiModelID).WillReturnRows(mockRows)
	foundAIModel, err := GetAIModel(mock, aiModelID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAIModel", err)
	}
	if cmp.Equal(*foundAIModel, targetData) == false {
		t.Errorf("Expected AIModel From Method GetAIModel: %v is different from actual %v", foundAIModel, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	aiModelID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WithArgs(aiModelID).WillReturnRows(noRows)
	foundAIModel, err := GetAIModel(mock, &aiModelID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAIModel", err)
	}
	if foundAIModel != nil {
		t.Errorf("Expected AIModel From Method GetAIModel: to be empty but got this: %v", foundAIModel)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelSqlErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	aiModelID := -1
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WithArgs(aiModelID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAIModel, err := GetAIModel(mock, &aiModelID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModel", err)
	}
	if foundAIModel != nil {
		t.Errorf("Expected AIModel From Method GetAIModel: to be empty but got this: %v", foundAIModel)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	aiModelID := -1
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WithArgs(aiModelID).WillReturnRows(differentModelRows)
	foundAIModel, err := GetAIModel(mock, &aiModelID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModel", err)
	}
	if foundAIModel != nil {
		t.Errorf("Expected AIModel From Method GetAIModel: to be empty but got this: %v", foundAIModel)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAIModel(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	aiModelID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM ai_models").WithArgs(*aiModelID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAIModel(mock, aiModelID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAIModel", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAIModelOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	aiModelID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveAIModel(mock, &aiModelID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveAIModelOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	aiModelID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM ai_models").WithArgs(aiModelID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveAIModel(mock, &aiModelID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	ids := []int{1, 2}
	mockRows := AddAIModelToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnRows(mockRows)
	foundAIModelList, err := GetAIModelList(mock, ids)
	if err != nil {
		t.Fatalf("an error '%s' in GetAIModelList", err)
	}
	for i, aiModelAIModel := range dataList {
		if cmp.Equal(aiModelAIModel, foundAIModelList[i]) == false {
			t.Errorf("Expected AIModel From Method GetAIModelList: %v is different from actual %v", aiModelAIModel, foundAIModelList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAIModelList, err := GetAIModelList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModelList", err)
	}
	if len(foundAIModelList) != 0 {
		t.Errorf("Expected From Method GetAIModelList: to be empty but got this: %v", foundAIModelList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	ids := []int{1, 2}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnRows(differentModelRows)
	foundAIModel, err := GetAIModelList(mock, ids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModelList", err)
	}
	if foundAIModel != nil {
		t.Errorf("Expected AIModel From Method GetAIModelList: to be empty but got this: %v", foundAIModel)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateAIModel(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE ai_models").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.OllamaName,    //6
		targetData.ParamsSize,    //7
		targetData.QuantizedSize, //8
		targetData.BaseModelID,   //9
		targetData.UpdatedBy,     //10
		targetData.ID,            //11
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAIModel(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAIModel", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAIModelOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateAIModel(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAIModelOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateAIModel(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateAIModelOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE ai_models").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.OllamaName,    //6
		targetData.ParamsSize,    //7
		targetData.QuantizedSize, //8
		targetData.BaseModelID,   //9
		targetData.UpdatedBy,     //10
		targetData.ID,            //11
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateAIModel(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModel(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.UUID = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO ai_models").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.OllamaName,    //6
		targetData.ParamsSize,    //7
		targetData.QuantizedSize, //8
		targetData.BaseModelID,   //9
		targetData.CreatedBy,     //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"ai_model_id"}).AddRow(1))
	mock.ExpectCommit()
	aiModelID, err := InsertAIModel(mock, &targetData)
	if aiModelID < 0 {
		t.Fatalf("aiModelID should not be negative ID: %d", aiModelID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertAIModel", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModelOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, err = InsertAIModel(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModelOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO ai_models").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.OllamaName,    //6
		targetData.ParamsSize,    //7
		targetData.QuantizedSize, //8
		targetData.BaseModelID,   //9
		targetData.CreatedBy,     //10
		targetData.CreatedBy,     //11
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	aiModelID, err := InsertAIModel(mock, &targetData)
	if aiModelID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error aiModelID: %d", aiModelID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModelOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO ai_models").WithArgs(
		targetData.Name,          //1
		targetData.AlternateName, //2
		targetData.URL,           //3
		targetData.Ticker,        //4
		targetData.Description,   //5
		targetData.OllamaName,    //6
		targetData.ParamsSize,    //7
		targetData.QuantizedSize, //8
		targetData.BaseModelID,   //9
		targetData.CreatedBy,     //10
		targetData.CreatedBy,     //11
	).WillReturnRows(pgxmock.NewRows([]string{"ai_model_id"}).AddRow(-1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	aiModelID, err := InsertAIModel(mock, &targetData)
	if aiModelID >= 0 {
		t.Fatalf("Expecting -1 for aiModelID because of error aiModelID: %d", aiModelID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModels(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"ai_models"}, DBColumnsInsertAIModels)
	targetData := TestAllData
	err = InsertAIModels(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAIModels", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertAIModelsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"ai_models"}, DBColumnsInsertAIModels).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertAIModels(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddAIModelToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnRows(mockRows)
	foundAIModelList, err := GetAIModelListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAIModelListByPagination", err)
	}
	for i, aiModelData := range dataList {
		if cmp.Equal(aiModelData, foundAIModelList[i]) == false {
			t.Errorf("Expected aiModelData From Method GetAIModelListByPagination: %v is different from actual %v", aiModelData, foundAIModelList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundAIModelList, err := GetAIModelListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModelListByPagination", err)
	}
	if len(foundAIModelList) != 0 {
		t.Errorf("Expected From Method GetAIModelListByPagination: to be empty but got this: %v", foundAIModelList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAIModelListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"ticker = 'COIN'", "name = 'coingecko'"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM ai_models").WillReturnRows(differentModelRows)
	foundAIModelList, err := GetAIModelListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAIModelListByPagination", err)
	}
	if len(foundAIModelList) != 0 {
		t.Errorf("Expected From Method GetAIModelListByPagination: to be empty but got this: %v", foundAIModelList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalAIModelsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM ai_models").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalAIModelsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalAIModelsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalAIModelsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM ai_models").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalAIModelsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalAIModelsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalAIModelsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
