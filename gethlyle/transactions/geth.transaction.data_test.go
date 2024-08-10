package gethlyletransactions

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                             //1
	"uuid",                           //2
	"chain_id",                       //3
	"exchange_id",                    //4
	"block_number",                   //5
	"index_number",                   //6
	"txn_date",                       //7
	"txn_hash",                       //8
	"from_address",                   //9
	"from_address_id",                //10
	"to_address",                     //11
	"to_address_id",                  //12
	"interacted_contract_address",    //13
	"interacted_contract_address_id", //14
	"native_asset_id",                //15
	"geth_process_job_id",            //16
	"value",                          //17
	"geth_transction_input_id",       //18
	"status_id",                      //19
	"description",                    //20
	"created_by",                     //21
	"created_at",                     //22
	"updated_by",                     //23
	"updated_at",                     //24
}
var DBColumnsInsertGethTransactions = []string{
	"uuid",                           //1
	"chain_id",                       //2
	"exchange_id",                    //3
	"block_number",                   //4
	"index_number",                   //5
	"txn_date",                       //6
	"txn_hash",                       //7
	"from_address",                   //8
	"from_address_id",                //9
	"to_address",                     //10
	"to_address_id",                  //11
	"interacted_contract_address",    //12
	"interacted_contract_address_id", //13
	"native_asset_id",                //14
	"geth_process_job_id",            //15
	"value",                          //16
	"geth_transction_input_id",       //17
	"status_id",                      //18
	"description",                    //19
	"created_by",                     //20
	"created_at",                     //21
	"updated_by",                     //22
	"updated_at",                     //23
}

var TestData1 = GethTransaction{
	ID:                          utils.Ptr[int](1),
	UUID:                        "01ef85e8-2c26-441e-8c7f-71d79518ad72",
	ChainID:                     utils.Ptr[int](1),
	ExchangeID:                  utils.Ptr[int](2),
	BlockNumber:                 utils.Ptr[uint64](20264466),
	IndexNumber:                 utils.Ptr[uint](1),
	TxnDate:                     utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	TxnHash:                     "0x6c695fdffb5063c3cb7ea3aef902cd1dbe9135cf14bdd5995c4a9698191fcc7c",
	FromAddress:                 "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326",
	FromAddressID:               utils.Ptr[int](1),
	ToAddress:                   "0xA8C62111e4652b07110A0FC81816303c42632f64",
	ToAddressID:                 utils.Ptr[int](2),
	InteractedContractAddress:   "",
	InteractedContractAddressID: nil,
	NativeAssetID:               utils.Ptr[int](1),
	GethProcessJobID:            utils.Ptr[int](10),
	Value:                       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(0.01)),
	GethTransctionInputId:       utils.Ptr[int](1),
	StatusID:                    utils.Ptr[int](1),
	Description:                 "",
	CreatedBy:                   "SYSTEM",
	CreatedAt:                   utils.SampleCreatedAtTime,
	UpdatedBy:                   "SYSTEM",
	UpdatedAt:                   utils.SampleCreatedAtTime,
}

var TestData2 = GethTransaction{
	ID:                          utils.Ptr[int](2),
	UUID:                        "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",
	ChainID:                     utils.Ptr[int](1),
	ExchangeID:                  utils.Ptr[int](2),
	BlockNumber:                 utils.Ptr[uint64](20272060),
	IndexNumber:                 utils.Ptr[uint](0),
	TxnDate:                     utils.Ptr[time.Time](utils.SampleCreatedAtTime),
	TxnHash:                     "0xfefca32d87fc4175d203646359fdb00643b57b8948fb3777dffa79d138f0b2c5",
	FromAddress:                 "0x93d8A622Fe3CC8477BBd22E205F3951f67FD64dE",
	FromAddressID:               utils.Ptr[int](4),
	ToAddress:                   "0x38C11FBaE0cf57B55a00951fBA9e6D1FDB4805DB",
	ToAddressID:                 utils.Ptr[int](5),
	InteractedContractAddress:   "",
	InteractedContractAddressID: nil,
	NativeAssetID:               utils.Ptr[int](1),
	GethProcessJobID:            utils.Ptr[int](2),
	Value:                       utils.Ptr[decimal.Decimal](decimal.NewFromFloat(11.2)),
	GethTransctionInputId:       utils.Ptr[int](555),
	StatusID:                    utils.Ptr[int](4),
	Description:                 "",
	CreatedBy:                   "SYSTEM",
	CreatedAt:                   utils.SampleCreatedAtTime,
	UpdatedBy:                   "SYSTEM",
	UpdatedAt:                   utils.SampleCreatedAtTime,
}
var TestAllData = []GethTransaction{TestData1, TestData2}

func AddGethTransactionToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTransaction) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,                          //1
			data.UUID,                        //2
			data.ChainID,                     //3
			data.ExchangeID,                  //4
			data.BlockNumber,                 //5
			data.IndexNumber,                 //6
			data.TxnDate,                     //7
			data.TxnHash,                     //8
			data.FromAddress,                 //9
			data.FromAddressID,               //10
			data.ToAddress,                   //11
			data.ToAddressID,                 //12
			data.InteractedContractAddress,   //13
			data.InteractedContractAddressID, //14
			data.NativeAssetID,               //15
			data.GethProcessJobID,            //16
			data.Value,                       //17
			data.GethTransctionInputId,       //18
			data.StatusID,                    //19
			data.Description,                 //20
			data.CreatedBy,                   //21
			data.CreatedAt,                   //22
			data.UpdatedBy,                   //23
			data.UpdatedAt,                   //24
		)
	}
	return rows
}

func TestGetGethTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethTransaction{targetData}
	gethTransactionID := targetData.ID
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(*gethTransactionID).WillReturnRows(mockRows)
	foundGethTransaction, err := GetGethTransaction(mock, gethTransactionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransaction", err)
	}
	if cmp.Equal(*foundGethTransaction, targetData) == false {
		t.Errorf("Expected GethTransaction From Method GetGethTransaction: %v is different from actual %v", foundGethTransaction, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(gethTransactionID).WillReturnRows(noRows)
	foundGethTransaction, err := GetGethTransaction(mock, &gethTransactionID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransaction", err)
	}
	if foundGethTransaction != nil {
		t.Errorf("Expected GethTransaction From Method GetGethTransaction: to be empty but got this: %v", foundGethTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(gethTransactionID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransaction, err := GetGethTransaction(mock, &gethTransactionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransaction", err)
	}
	if foundGethTransaction != nil {
		t.Errorf("Expected GethTransaction From Method GetGethTransaction: to be empty but got this: %v", foundGethTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionID := -1

	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(gethTransactionID).WillReturnRows(differentModelRows)
	foundGethTransaction, err := GetGethTransaction(mock, &gethTransactionID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransaction", err)
	}
	if foundGethTransaction != nil {
		t.Errorf("Expected foundGethTransaction From Method GetGethTransaction: to be empty but got this: %v", foundGethTransaction)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromToAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransaction{TestData1, TestData2}
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	fromToAddressID := TestData1.FromAddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(*fromToAddressID).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionByFromToAddress(mock, fromToAddressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionByFromToAddress", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionByFromToAddress: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromToAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromToAddressID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(fromToAddressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionByFromToAddress(mock, &fromToAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionByFromToAddress", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionByFromToAddress: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromToAddressForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromToAddressID := -999
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(fromToAddressID).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionByFromToAddress(mock, &fromToAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionByFromToAddress", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionByFromToAddress: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromAddressAndBeforeBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransaction{TestData1, TestData2}
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	fromToAddressID := TestData1.FromAddressID
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(*fromToAddressID, *blockNumber).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionByFromAddressAndBeforeBlockNumber(mock, fromToAddressID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionByFromAddressAndBeforeBlockNumber", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionByFromAddressAndBeforeBlockNumber: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromAddressAndBeforeBlockNumberForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromToAddressID := -999
	blockNumber := uint64(20272060)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(fromToAddressID, blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionByFromAddressAndBeforeBlockNumber(mock, &fromToAddressID, &blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionByFromAddressAndBeforeBlockNumber", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionByFromAddressAndBeforeBlockNumber: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionByFromAddressAndBeforeBlockNumberForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	fromToAddressID := -999
	blockNumber := uint64(20272060)
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(fromToAddressID, blockNumber).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionByFromAddressAndBeforeBlockNumber(mock, &fromToAddressID, &blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionByFromAddressAndBeforeBlockNumber", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionByFromAddressAndBeforeBlockNumber: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHash(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	txnHash := TestData1.TxnHash
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(txnHash).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionsByTxnHash(mock, txnHash)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionsByTxnHash", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionsByTxnHash: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHashForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHash := "0x-invalid"
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(txnHash).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionsByTxnHash(mock, txnHash)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByTxnHash", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionsByTxnHash: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHashForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHash := "0x-invalid"
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(txnHash).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionsByTxnHash(mock, txnHash)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByTxnHash", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionsByTxnHash: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHashes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(txnHashes)).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionsByTxnHashes(mock, txnHashes)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionsByTxnHashes", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionsByTxnHashes: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHashesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	txnHashes := []string{"0x-invalid-1", "0x-invalid-2"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(txnHashes)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionsByTxnHashes(mock, txnHashes)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByTxnHashes", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionsByTxnHashes: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByTxnHashesForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHashes := []string{"0x-invalid-1", "0x-invalid-2"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(txnHashes)).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionsByTxnHashes(mock, txnHashes)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByTxnHashes", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionsByTxnHashes: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByUUIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	uuids := []string{TestData1.UUID, TestData2.UUID}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(uuids)).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionsByUUIDs(mock, uuids)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionsByUUIDs", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionsByUUIDs: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByUUIDsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	uuids := []string{"uuid-invalid-1", "uuid-invalid-2"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(uuids)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByUUIDs", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionsByUUIDs: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionsByUUIDsForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	uuids := []string{"uuid-invalid-1", "uuid-invalid-2"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WithArgs(pq.Array(uuids)).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionsByUUIDs(mock, uuids)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionsByUUIDs", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionsByUUIDs: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethTransactionID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(*gethTransactionID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransaction(mock, gethTransactionID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethTransaction(mock, &gethTransactionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransactionID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(gethTransactionID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransaction(mock, &gethTransactionID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionsFromChainIDAndStartBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	chainID := targetData.ChainID
	startBlockNumber := targetData.BlockNumber
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(*chainID, *startBlockNumber).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransactionsFromChainIDAndStartBlockNumber(mock, chainID, startBlockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransactionsFromChainIDAndStartBlockNumber", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionsFromChainIDAndStartBlockNumberOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	startBlockNumber := uint64(10000)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethTransactionsFromChainIDAndStartBlockNumber(mock, &chainID, &startBlockNumber)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionsFromChainIDAndStartBlockNumberOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	startBlockNumber := uint64(10000)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(chainID, startBlockNumber).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransactionsFromChainIDAndStartBlockNumber(mock, &chainID, &startBlockNumber)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionsFromChainID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	chainID := targetData.ChainID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(*chainID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransactionsFromChainID(mock, chainID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransactionsFromChainID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransactionsFromChainIDOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = RemoveGethTransactionsFromChainID(mock, &chainID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestRemoveGethTransactionsFromChainIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	chainID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transactions").WithArgs(chainID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransactionsFromChainID(mock, &chainID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddGethTransactionToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionList", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetGethTransactionList: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionListForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionList", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetGethTransactionList: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionList", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionList: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transactions").WithArgs(
		targetData.ChainID,                     //1
		targetData.ExchangeID,                  //2
		targetData.BlockNumber,                 //3
		targetData.IndexNumber,                 //4
		targetData.TxnDate,                     //5
		targetData.TxnHash,                     //6
		targetData.FromAddress,                 //7
		targetData.FromAddressID,               //8
		targetData.ToAddress,                   //9
		targetData.ToAddressID,                 //10
		targetData.InteractedContractAddress,   //11
		targetData.InteractedContractAddressID, //12
		targetData.NativeAssetID,               //13
		targetData.GethProcessJobID,            //14
		targetData.Value,                       //15
		targetData.GethTransctionInputId,       //16
		targetData.StatusID,                    //17
		targetData.Description,                 //18
		targetData.UpdatedBy,                   //19
		targetData.ID,                          //20
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransaction(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransactionOnFailureAtParameter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = nil
	err = UpdateGethTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transactions").WithArgs(
		targetData.ChainID,                     //1
		targetData.ExchangeID,                  //2
		targetData.BlockNumber,                 //3
		targetData.IndexNumber,                 //4
		targetData.TxnDate,                     //5
		targetData.TxnHash,                     //6
		targetData.FromAddress,                 //7
		targetData.FromAddressID,               //8
		targetData.ToAddress,                   //9
		targetData.ToAddressID,                 //10
		targetData.InteractedContractAddress,   //11
		targetData.InteractedContractAddressID, //12
		targetData.NativeAssetID,               //13
		targetData.GethProcessJobID,            //14
		targetData.Value,                       //15
		targetData.GethTransctionInputId,       //16
		targetData.StatusID,                    //17
		targetData.Description,                 //18
		targetData.UpdatedBy,                   //19
		targetData.ID,                          //20
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transactions").WithArgs(
		targetData.ChainID,                     //1
		targetData.ExchangeID,                  //2
		targetData.BlockNumber,                 //3
		targetData.IndexNumber,                 //4
		targetData.TxnDate,                     //5
		targetData.TxnHash,                     //6
		targetData.FromAddress,                 //7
		targetData.FromAddressID,               //8
		targetData.ToAddress,                   //9
		targetData.ToAddressID,                 //10
		targetData.InteractedContractAddress,   //11
		targetData.InteractedContractAddressID, //12
		targetData.NativeAssetID,               //13
		targetData.GethProcessJobID,            //14
		targetData.Value,                       //15
		targetData.GethTransctionInputId,       //16
		targetData.StatusID,                    //17
		targetData.Description,                 //18
		targetData.CreatedBy,                   //19
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	gethTransactionID, newUUID, err := InsertGethTransaction(mock, &targetData)
	if gethTransactionID < 0 {
		t.Fatalf("gethTransactionID should not be negative ID: %d", gethTransactionID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransaction", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	_, _, err = InsertGethTransaction(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestInsertGethTransactionOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.BlockNumber = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transactions").WithArgs(
		targetData.ChainID,                     //1
		targetData.ExchangeID,                  //2
		targetData.BlockNumber,                 //3
		targetData.IndexNumber,                 //4
		targetData.TxnDate,                     //5
		targetData.TxnHash,                     //6
		targetData.FromAddress,                 //7
		targetData.FromAddressID,               //8
		targetData.ToAddress,                   //9
		targetData.ToAddressID,                 //10
		targetData.InteractedContractAddress,   //11
		targetData.InteractedContractAddressID, //12
		targetData.NativeAssetID,               //13
		targetData.GethProcessJobID,            //14
		targetData.Value,                       //15
		targetData.GethTransctionInputId,       //16
		targetData.StatusID,                    //17
		targetData.Description,                 //18
		targetData.CreatedBy,                   //19
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTransactionID, newUUID, err := InsertGethTransaction(mock, &targetData)
	if gethTransactionID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTransactionID: %d", gethTransactionID)
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

func TestInsertGethTransactionOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transactions").WithArgs(
		targetData.ChainID,                     //1
		targetData.ExchangeID,                  //2
		targetData.BlockNumber,                 //3
		targetData.IndexNumber,                 //4
		targetData.TxnDate,                     //5
		targetData.TxnHash,                     //6
		targetData.FromAddress,                 //7
		targetData.FromAddressID,               //8
		targetData.ToAddress,                   //9
		targetData.ToAddressID,                 //10
		targetData.InteractedContractAddress,   //11
		targetData.InteractedContractAddressID, //12
		targetData.NativeAssetID,               //13
		targetData.GethProcessJobID,            //14
		targetData.Value,                       //15
		targetData.GethTransctionInputId,       //16
		targetData.StatusID,                    //17
		targetData.Description,                 //18
		targetData.CreatedBy,                   //19
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTransactionID, newUUID, err := InsertGethTransaction(mock, &targetData)
	if gethTransactionID >= 0 {
		t.Fatalf("Expecting -1 for gethTransactionID because of error gethTransactionID: %d", gethTransactionID)
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

func TestInsertGethTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transactions"}, DBColumnsInsertGethTransactions)
	targetData := TestAllData
	err = InsertGethTransactions(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransactions", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransactionsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transactions"}, DBColumnsInsertGethTransactions).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethTransactions(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransactionAddresses(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transactions").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transactions").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransactionAddresses(mock)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransactionAddresses", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionAddressesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransactionAddresses(mock)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransactionAddressesOnFailure1st(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transactions").WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTransactionAddresses(mock)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransactionAddressesOnFailure2nd(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transactions").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transactions").WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTransactionAddresses(mock)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromTransactions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []string{TestData1.FromAddress, TestData2.FromAddress}
	mockRows := mock.NewRows([]string{"address"}).AddRow(TestData1.FromAddress).AddRow(TestData2.FromAddress)
	mock.ExpectQuery("^WITH sender_table as ").WillReturnRows(mockRows)
	foundNullAddresses, err := GetNullAddressStrsFromTransactions(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetNullAddressStrsFromTransactions", err)
	}
	for i, nullAddress := range foundNullAddresses {
		if cmp.Equal(nullAddress, dataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetNullAddressStrsFromTransactions: %v is different from actual %v", nullAddress, dataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromTransactionsForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^WITH sender_table as ").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetNullAddressStrsFromTransactions(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetNullAddressStrsFromTransactions", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetNullAddressStrsFromTransactions: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	_start := 1
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = 1", "exchange_id=1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnRows(mockRows)
	foundGethTransactionList, err := GetGethTransactionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransactionListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethTransactionList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethTransactionListByPagination: %v is different from actual %v", sourceData, foundGethTransactionList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionListByPaginationForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = -1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetGethTransactionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionListByPagination", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetGethTransactionListByPagination: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransactionListByPaginationForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = 1"}
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transactions").WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetGethTransactionListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransactionListByPagination", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected From Method GetGethTransactionListByPagination: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTotalTransactionsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	numOfChainsExpected := 10
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transactions").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetTotalTransactionsCount", err)
	}
	if *numOfChains != numOfChainsExpected {
		t.Errorf("Expected Chain From Method GetTotalTransactionsCount: %d is different from actual %d", numOfChainsExpected, *numOfChains)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transactions").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	numOfChains, err := GetTotalTransactionsCount(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTotalTransactionsCount", err)
	}
	if numOfChains != nil {
		t.Errorf("Expected numOfChains From Method GetTotalTransactionsCount to be empty but got this: %v", numOfChains)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	minerID := 1
	fromAddress := "Test-From-Address"
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddress(mock, &minerID, fromAddress)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddress", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetAllGethTransactionsByMinerIDAndFromAddress: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddress(mock, &minerID, fromAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddress", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetAllGethTransactionsByMinerIDAndFromAddress: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddress(mock, &minerID, fromAddress)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddress", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetAllGethTransactionsByMinerIDAndFromAddress: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressToDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	minerID := 1
	fromAddress := "Test-From-Address"
	toDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, toDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressToDate(mock, &minerID, fromAddress, &toDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressToDate", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetAllGethTransactionsByMinerIDAndFromAddressToDate: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressToDateForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	toDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, toDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressToDate(mock, &minerID, fromAddress, &toDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressToDate", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetAllGethTransactionsByMinerIDAndFromAddressToDate: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressToDateForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	toDate := utils.SampleCreatedAtTime
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, toDate.Format(utils.LayoutPostgres)).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressToDate(mock, &minerID, fromAddress, &toDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressToDate", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetAllGethTransactionsByMinerIDAndFromAddressToDate: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressFromToDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransactionToMockRows(mock, dataList)
	minerID := 1
	fromAddress := "Test-From-Address"
	fromDate := utils.SampleCreatedAtTime
	toDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, fromDate.Format(utils.LayoutPostgres), toDate.Format(utils.LayoutPostgres)).WillReturnRows(mockRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressFromToDate(mock, &minerID, fromAddress, &fromDate, &toDate)
	if err != nil {
		t.Fatalf("an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressFromToDate", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransaction := range foundGethTransactionList {
		if cmp.Equal(foundGethTransaction, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransaction From Method GetAllGethTransactionsByMinerIDAndFromAddressFromToDate: %v is different from actual %v", foundGethTransaction, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressFromToDateForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	fromDate := utils.SampleCreatedAtTime
	toDate := utils.SampleCreatedAtTime
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, fromDate.Format(utils.LayoutPostgres), toDate.Format(utils.LayoutPostgres)).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressFromToDate(mock, &minerID, fromAddress, &fromDate, &toDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressFromToDate", err)
	}
	if len(foundGethTransactionList) != 0 {
		t.Errorf("Expected From Method GetAllGethTransactionsByMinerIDAndFromAddressFromToDate: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetAllGethTransactionsByMinerIDAndFromAddressFromToDateForCollectRowsErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	minerID := -1
	fromAddress := "Test-From-Address"
	fromDate := utils.SampleCreatedAtTime
	toDate := utils.SampleCreatedAtTime
	differentModelRows := mock.NewRows([]string{"diff_model_id"}).AddRow(1)
	mock.ExpectQuery("^SELECT (.+) FROM geth_miners_transactions gmt LEFT JOIN geth_transactions gt ON gmt.transaction_id = gt.id").WithArgs(minerID, fromAddress, fromDate.Format(utils.LayoutPostgres), toDate.Format(utils.LayoutPostgres)).WillReturnRows(differentModelRows)
	foundGethTransactionList, err := GetAllGethTransactionsByMinerIDAndFromAddressFromToDate(mock, &minerID, fromAddress, &fromDate, &toDate)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetAllGethTransactionsByMinerIDAndFromAddressFromToDate", err)
	}
	if foundGethTransactionList != nil {
		t.Errorf("Expected foundGethTransactionList From Method GetAllGethTransactionsByMinerIDAndFromAddressFromToDate: to be empty but got this: %v", foundGethTransactionList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
