package gethlyleminers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",                    //1
	"uuid",                  //2
	"name",                  //3
	"alternate_name",        //4
	"chain_id",              //5
	"exchange_id",           //6
	"starting_block_number", //7
	"created_txn_hash",      //8
	"last_block_number",     //9
	"contract_address",      //10
	"contract_address_id",   //11
	"developer_address",     //12
	"developer_address_id",  //13
	"mining_asset_id",       //14
	"description",           //15
	"created_by",            //16
	"created_at",            //17
	"updated_by",            //18
	"updated_at",            //19
}
var DBColumnsInsertGethMiners = []string{
	"uuid",                  //1
	"name",                  //2
	"alternate_name",        //3
	"chain_id",              //4
	"exchange_id",           //5
	"starting_block_number", //6
	"created_txn_hash",      //7
	"last_block_number",     //8
	"contract_address",      //9
	"contract_address_id",   //10
	"developer_address",     //11
	"developer_address_id",  //12
	"mining_asset_id",       //13
	"description",           //14
	"created_by",            //15
	"created_at",            //16
	"updated_by",            //17
	"updated_at",            //18
}

var TestData1 = GethMiner{
	ID:            utils.Ptr[int](1),
	UUID:          "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	Name:          "Meow Miner",
	AlternateName: "Meow",

	ChainID:             utils.Ptr[int](13),
	ExchangeID:          utils.Ptr[int](4),
	StartingBlockNumber: utils.Ptr[int](42830364),
	CreatedTxnHash:      "0x19c14e99d55adc44750791d7532b98a577cc8877ff04908a4eb45b58bfea97f1",
	LastBlockNumber:     utils.Ptr[uint64](44185364),
	ContractAddress:     "0xc0F9a97E46Fb0f80aE39981759eAB4a61eE36459",
	ContractAddressID:   utils.Ptr[int](1),
	DeveloperAddress:    "0xBD1C7f2A06aC1C7cec56A63B3c1c2CeE8C50e92d",
	DeveloperAddressID:  utils.Ptr[int](1),
	MiningAssetID:       utils.Ptr[int](15797),
	Description:         "https://meowminer.com/",
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
}

var TestData2 = GethMiner{
	ID:                  utils.Ptr[int](2),
	UUID:                "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	Name:                "Print The PEPE",
	AlternateName:       "Print The PEPE",
	ChainID:             utils.Ptr[int](1),
	ExchangeID:          utils.Ptr[int](2),
	StartingBlockNumber: utils.Ptr[int](42830333),
	CreatedTxnHash:      "0x19c14e99d55adc44750791d7532b98a577cc8877ff04908a4eb45b58bfea9123",
	LastBlockNumber:     utils.Ptr[uint64](44185364),
	ContractAddress:     "0xc0F9a97E46Fb0f80aE39981759eAB4a61eE36123",
	ContractAddressID:   nil,
	DeveloperAddress:    "0xBD1C7f2A06aC1C7cec56A63B3c1c2CeE8C50e123",
	DeveloperAddressID:  nil,
	MiningAssetID:       utils.Ptr[int](15733),
	Description:         "https://printthepepe.com/",
	CreatedBy:           "SYSTEM",
	CreatedAt:           utils.SampleCreatedAtTime,
	UpdatedBy:           "SYSTEM",
	UpdatedAt:           utils.SampleCreatedAtTime,
}
var TestAllData = []GethMiner{TestData1, TestData2}

func AddGethMinerToMockRows(mock pgxmock.PgxPoolIface, dataList []GethMiner) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                  //1
			data.UUID,                //2
			data.Name,                //3
			data.AlternateName,       //4
			data.ChainID,             //5
			data.ExchangeID,          //6
			data.StartingBlockNumber, //7
			data.CreatedTxnHash,      //8
			data.LastBlockNumber,     //9
			data.ContractAddress,     //10
			data.ContractAddressID,   //11
			data.DeveloperAddress,    //12
			data.DeveloperAddressID,  //13
			data.MiningAssetID,       //14
			data.Description,         //15
			data.CreatedBy,           //16
			data.CreatedAt,           //17
			data.UpdatedBy,           //18
			data.UpdatedAt,           //19
		)
	}
	return rows
}

func TestGetGethMiner(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethMiner{targetData}
	gethMinerID := targetData.ID
	mockRows := AddGethMinerToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WithArgs(*gethMinerID).WillReturnRows(mockRows)
	foundMarketData, err := GetGethMiner(mock, gethMinerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMiner", err)
	}
	if cmp.Equal(*foundMarketData, targetData) == false {
		t.Errorf("Expected GethMiner From Method GetGethMiner: %v is different from actual %v", foundMarketData, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WithArgs(gethMinerID).WillReturnRows(noRows)
	foundMarketData, err := GetGethMiner(mock, &gethMinerID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMiner", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMiner From Method GetGethMiner: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WithArgs(gethMinerID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketData, err := GetGethMiner(mock, &gethMinerID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMiner", err)
	}
	if foundMarketData != nil {
		t.Errorf("Expected GethMiner From Method GetGethMiner: to be empty but got this: %v", foundMarketData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMiner(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethMinerID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners").WithArgs(*gethMinerID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethMiner(mock, gethMinerID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethMiner", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethMinerOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethMinerID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_miners").WithArgs(gethMinerID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethMiner(mock, &gethMinerID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMiner{TestData1, TestData2}
	mockRows := AddGethMinerToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WillReturnRows(mockRows)
	foundMarketDataList, err := GetGethMinerList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerList", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMiner From Method GetGethMinerList: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetGethMinerList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerList", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMinerList: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerListByMiningAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethMiner{TestData1}
	mockRows := AddGethMinerToMockRows(mock, dataList)
	miningAssetID := TestData1.MiningAssetID
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WithArgs(*miningAssetID).WillReturnRows(mockRows)
	foundMarketDataList, err := GetGethMinerListByMiningAssetId(mock, miningAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerListByMiningAssetId", err)
	}
	testMarketDataList := TestAllData
	for i, foundMarketData := range foundMarketDataList {
		if cmp.Equal(foundMarketData, testMarketDataList[i]) == false {
			t.Errorf("Expected GethMiner From Method GetGethMinerListByMiningAssetId: %v is different from actual %v", foundMarketData, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerListByMiningAssetIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	miningAssetID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WithArgs(miningAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundMarketDataList, err := GetGethMinerListByMiningAssetId(mock, &miningAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerListByMiningAssetId", err)
	}
	if len(foundMarketDataList) != 0 {
		t.Errorf("Expected From Method GetGethMinerListByMiningAssetId: to be empty but got this: %v", foundMarketDataList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethMiner(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.ChainID,             //3
		targetData.ExchangeID,          //4
		targetData.StartingBlockNumber, //5
		targetData.CreatedTxnHash,      //6
		targetData.LastBlockNumber,     //7
		targetData.ContractAddress,     //8
		targetData.ContractAddressID,   //9
		targetData.DeveloperAddress,    //10
		targetData.DeveloperAddressID,  //11
		targetData.MiningAssetID,       //12
		targetData.Description,         //13
		targetData.UpdatedBy,           //14
		targetData.ID,                  //15
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethMiner(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethMiner", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethMiner(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethMinerOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_miners").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.ChainID,             //3
		targetData.ExchangeID,          //4
		targetData.StartingBlockNumber, //5
		targetData.CreatedTxnHash,      //6
		targetData.LastBlockNumber,     //7
		targetData.ContractAddress,     //8
		targetData.ContractAddressID,   //9
		targetData.DeveloperAddress,    //10
		targetData.DeveloperAddressID,  //11
		targetData.MiningAssetID,       //12
		targetData.Description,         //13
		targetData.UpdatedBy,           //14
		targetData.ID,                  //15
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethMiner(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMiner(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.Name = "New Name"
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.ChainID,             //3
		targetData.ExchangeID,          //4
		targetData.StartingBlockNumber, //5
		targetData.CreatedTxnHash,      //6
		targetData.LastBlockNumber,     //7
		targetData.ContractAddress,     //8
		targetData.ContractAddressID,   //9
		targetData.DeveloperAddress,    //10
		targetData.DeveloperAddressID,  //11
		targetData.MiningAssetID,       //12
		targetData.Description,         //13
		targetData.CreatedBy,           //14
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	gethMinerID, newUUID, err := InsertGethMiner(mock, &targetData)
	if gethMinerID < 0 {
		t.Fatalf("gethMinerID should not be negative ID: %d", gethMinerID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMiner", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.ChainID,             //3
		targetData.ExchangeID,          //4
		targetData.StartingBlockNumber, //5
		targetData.CreatedTxnHash,      //6
		targetData.LastBlockNumber,     //7
		targetData.ContractAddress,     //8
		targetData.ContractAddressID,   //9
		targetData.DeveloperAddress,    //10
		targetData.DeveloperAddressID,  //11
		targetData.MiningAssetID,       //12
		targetData.Description,         //13
		targetData.CreatedBy,           //14
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, newUUID, err := InsertGethMiner(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethMinerID: %d", gethMinerID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinerOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.Name = ""
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_miners").WithArgs(
		targetData.Name,                //1
		targetData.AlternateName,       //2
		targetData.ChainID,             //3
		targetData.ExchangeID,          //4
		targetData.StartingBlockNumber, //5
		targetData.CreatedTxnHash,      //6
		targetData.LastBlockNumber,     //7
		targetData.ContractAddress,     //8
		targetData.ContractAddressID,   //9
		targetData.DeveloperAddress,    //10
		targetData.DeveloperAddressID,  //11
		targetData.MiningAssetID,       //12
		targetData.Description,         //13
		targetData.CreatedBy,           //14
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethMinerID, newUUID, err := InsertGethMiner(mock, &targetData)
	if gethMinerID >= 0 {
		t.Fatalf("Expecting -1 for gethMinerID because of error gethMinerID: %d", gethMinerID)
	}
	if newUUID != "" {
		t.Fatalf("on failure newUUID should be empty string: %s", newUUID)
	}
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMiners(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners"}, DBColumnsInsertGethMiners)
	targetData := TestAllData
	err = InsertGethMiners(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethMiners", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethMinersOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_miners"}, DBColumnsInsertGethMiners).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethMiners(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethMinerToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WillReturnRows(mockRows)
	foundChains, err := GetGethMinerListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethMinerListByPagination", err)
	}
	testChains := dataList
	for i, foundChain := range foundChains {
		if cmp.Equal(foundChain, testChains[i]) == false {
			t.Errorf("Expected Chain From Method GetGethMinerListByPagination: %v is different from actual %v", foundChain, testChains[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethMinerListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"import_type_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundChains, err := GetGethMinerListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethMinerListByPagination", err)
	}
	if len(foundChains) != 0 {
		t.Errorf("Expected From Method GetGethMinerListByPagination: to be empty but got this: %v", foundChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalGethMinersCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalGethMinersCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalGethMinersCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalGethMinersCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_miners").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalGethMinersCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalGethMinersCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalGethMinersCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
