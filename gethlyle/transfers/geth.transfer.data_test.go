package gethlyletransfers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	gethlyleaddresses "github.com/kfukue/lyle-labs-libraries/gethlyle/address"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
)

var DBColumns = []string{
	"id",                  //1
	"uuid",                //2
	"chain_id",            //3
	"token_address",       //4
	"token_address_id",    //5
	"asset_id",            //6
	"block_number",        //7
	"index_number",        //8
	"transfer_date",       //9
	"txn_hash",            //10
	"sender_address",      //11
	"sender_address_id",   //12
	"to_address",          //13
	"to_address_id",       //14
	"amount",              //15
	"description",         //16
	"created_by",          //17
	"created_at",          //18
	"updated_by",          //19
	"updated_at",          //20
	"geth_process_job_id", //21
	"topics_str",          //22
	"status_id",           //23
	"base_asset_id",       //24
	"transfer_type_id",    //25
}
var DBColumnsInsertGethTransfers = []string{
	"uuid",                //1
	"chain_id",            //2
	"token_address",       //3
	"token_address_id",    //4
	"asset_id",            //5
	"block_number",        //6
	"index_number",        //7
	"transfer_date",       //8
	"txn_hash",            //9
	"sender_address",      //10
	"sender_address_id",   //11
	"to_address",          //12
	"to_address_id",       //13
	"amount",              //14
	"description",         //15
	"created_by",          //16
	"created_at",          //17
	"updated_by",          //18
	"updated_at",          //19
	"geth_process_job_id", //20
	"topics_str",          //21
	"status_id",           //22
	"base_asset_id",       //23
	"transfer_type_id",    //24
}

var TestData1 = GethTransfer{
	ID:               utils.Ptr[int](1),                                                              //1
	UUID:             "01ef85e8-2c26-441e-8c7f-71d79518ad72",                                         //2
	ChainID:          utils.Ptr[int](1),                                                              //3
	TokenAddress:     "0xf819d9Cb1c2A819Fd991781A822dE3ca8607c3C9",                                   //4
	TokenAddressID:   utils.Ptr[int](1),                                                              //5
	AssetID:          utils.Ptr[int](1),                                                              //6
	BlockNumber:      utils.Ptr[uint64](17471544),                                                    //7
	IndexNumber:      utils.Ptr[uint](180),                                                           //8
	TransferDate:     utils.SampleCreatedAtTime,                                                      //9
	TxnHash:          "0xb16657fb1e468132f2cfcdb45bb53d1309097303d24b6446471bd802642a680f",           //10
	SenderAddress:    "0xaA6aeB29D2c55C3090E782123dADbE61eDcb2489",                                   //11
	SenderAddressID:  utils.Ptr[int](33911),                                                          //12
	ToAddress:        "0xD6cD56B69c2C345cdD98125ADE6d25E1161d1a2e",                                   //13
	ToAddressID:      utils.Ptr[int](34985),                                                          //14
	Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(137180999604867465199)),        //15
	Description:      "Imported by Geth Dex Analyzer",                                                //16
	CreatedBy:        "SYSTEM",                                                                       //17
	CreatedAt:        utils.SampleCreatedAtTime,                                                      //18
	UpdatedBy:        "SYSTEM",                                                                       //19
	UpdatedAt:        utils.SampleCreatedAtTime,                                                      //20
	GethProcessJobID: utils.Ptr[int](10),                                                             //21
	TopicsStr:        []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, //22
	StatusID:         utils.Ptr[int](52),                                                             //23
	BaseAssetID:      utils.Ptr[int](1),                                                              //24
	TransferTypeID:   utils.Ptr[int](2),                                                              //25

}

var TestData2 = GethTransfer{
	ID:               utils.Ptr[int](2),                                                              //1
	UUID:             "4f0d5402-7a7c-402d-a7fc-c56a02b13e03",                                         //2
	ChainID:          utils.Ptr[int](1),                                                              //3
	TokenAddress:     "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",                                   //4
	TokenAddressID:   utils.Ptr[int](2),                                                              //5
	AssetID:          utils.Ptr[int](3),                                                              //6
	BlockNumber:      utils.Ptr[uint64](17662045),                                                    //7
	IndexNumber:      utils.Ptr[uint](250),                                                           //8
	TransferDate:     utils.SampleCreatedAtTime,                                                      //9
	TxnHash:          "0x32d78b63a9a6e54f8078730b9d529cc2c6809a40f730e0bdf3339cb8704afb06",           //10
	SenderAddress:    "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",                                   //11
	SenderAddressID:  utils.Ptr[int](12237),                                                          //12
	ToAddress:        "0x479D7a7864e00668DA1ee507Fa1b925FB1b2AE19",                                   //13
	ToAddressID:      utils.Ptr[int](10535),                                                          //14
	Amount:           utils.Ptr[decimal.Decimal](decimal.NewFromFloat(31635761458482936)),            //15
	Description:      "Imported by Geth Dex Analyzer",                                                //16
	CreatedBy:        "SYSTEM",                                                                       //17
	CreatedAt:        utils.SampleCreatedAtTime,                                                      //18
	UpdatedBy:        "SYSTEM",                                                                       //19
	UpdatedAt:        utils.SampleCreatedAtTime,                                                      //20
	GethProcessJobID: utils.Ptr[int](10),                                                             //21
	TopicsStr:        []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, //22
	StatusID:         utils.Ptr[int](52),                                                             //23
	BaseAssetID:      utils.Ptr[int](1),                                                              //24
	TransferTypeID:   utils.Ptr[int](2),                                                              //25
}
var TestAllData = []GethTransfer{TestData1, TestData2}

func AddGethTransferToMockRows(mock pgxmock.PgxPoolIface, dataList []GethTransfer) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,               //1
			data.UUID,             //2
			data.ChainID,          //3
			data.TokenAddress,     //4
			data.TokenAddressID,   //5
			data.AssetID,          //6
			data.BlockNumber,      //7
			data.IndexNumber,      //8
			data.TransferDate,     //9
			data.TxnHash,          //10
			data.SenderAddress,    //11
			data.SenderAddressID,  //12
			data.ToAddress,        //13
			data.ToAddressID,      //14
			data.Amount,           //15
			data.Description,      //16
			data.CreatedBy,        //17
			data.CreatedAt,        //18
			data.UpdatedBy,        //19
			data.UpdatedAt,        //20
			data.GethProcessJobID, //21
			data.TopicsStr,        //22
			data.StatusID,         //23
			data.BaseAssetID,      //24
			data.TransferTypeID,   //25
		)
	}
	return rows
}

func TestGetGethTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethTransfer{targetData}
	gethTransferID := targetData.ID
	mockRows := AddGethTransferToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*gethTransferID).WillReturnRows(mockRows)
	foundGethTransfer, err := GetGethTransfer(mock, gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransfer", err)
	}
	if cmp.Equal(*foundGethTransfer, targetData) == false {
		t.Errorf("Expected GethTransfer From Method GetGethTransfer: %v is different from actual %v", foundGethTransfer, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransferID := 999
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(gethTransferID).WillReturnRows(noRows)
	foundGethTransfer, err := GetGethTransfer(mock, &gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransfer", err)
	}
	if foundGethTransfer != nil {
		t.Errorf("Expected GethTransfer From Method GetGethTransfer: to be empty but got this: %v", foundGethTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransferID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(gethTransferID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransfer, err := GetGethTransfer(mock, &gethTransferID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransfer", err)
	}
	if foundGethTransfer != nil {
		t.Errorf("Expected GethTransfer From Method GetGethTransfer: to be empty but got this: %v", foundGethTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByBlockChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	dataList := []GethTransfer{targetData}
	txnHash := targetData.TxnHash
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber

	mockRows := AddGethTransferToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(txnHash, *blockNumber, *indexNumber).WillReturnRows(mockRows)
	foundGethTransfer, err := GetGethTransferByBlockChain(mock, txnHash, blockNumber, indexNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByBlockChain", err)
	}
	if cmp.Equal(*foundGethTransfer, targetData) == false {
		t.Errorf("Expected GethTransfer From Method GetGethTransferByBlockChain: %v is different from actual %v", foundGethTransfer, targetData)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByBlockChainForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	txnHash := "non-existing-txn-hash"
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(txnHash, *blockNumber, *indexNumber).WillReturnRows(noRows)
	foundGethTransfer, err := GetGethTransferByBlockChain(mock, txnHash, blockNumber, indexNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByBlockChain", err)
	}
	if foundGethTransfer != nil {
		t.Errorf("Expected GethTransfer From Method GetGethTransferByBlockChain: to be empty but got this: %v", foundGethTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByBlockChainForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	txnHash := "invalid-txn-hash"
	blockNumber := targetData.BlockNumber
	indexNumber := targetData.IndexNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(txnHash, *blockNumber, *indexNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransfer, err := GetGethTransferByBlockChain(mock, txnHash, blockNumber, indexNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferByBlockChain", err)
	}
	if foundGethTransfer != nil {
		t.Errorf("Expected GethTransfer From Method GetGethTransferByBlockChain: to be empty but got this: %v", foundGethTransfer)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransfersTransactionHashByUserAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHashResults := []string{TestData1.TxnHash, TestData2.TxnHash}
	mockRows := mock.NewRows([]string{"txn_hash"}).AddRow(txnHashResults[0]).AddRow(txnHashResults[1])
	userAddressID := TestData1.ToAddressID
	assetID := TestData1.AssetID
	blockNumber := TestData1.BlockNumber

	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*userAddressID, *assetID, *blockNumber).WillReturnRows(mockRows)
	foundTxnHashes, err := GetTransfersTransactionHashByUserAddress(mock, userAddressID, assetID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetTransfersTransactionHashByUserAddress", err)
	}
	for i, sourceTxnHash := range txnHashResults {
		if cmp.Equal(sourceTxnHash, foundTxnHashes[i]) == false {
			t.Errorf("Expected foundTxnHashes From Method GetTransfersTransactionHashByUserAddress: %v is different from actual %v", sourceTxnHash, foundTxnHashes[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetTransfersTransactionHashByUserAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	userAddressID := -1
	assetID := -1
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(userAddressID, assetID, *blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetTransfersTransactionHashByUserAddress(mock, &userAddressID, &assetID, blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetTransfersTransactionHashByUserAddress", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetTransfersTransactionHashByUserAddress: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctAddressesFromAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethAddressResults := []gethlyleaddresses.GethAddress{gethlyleaddresses.TestData1, gethlyleaddresses.TestData2}
	mockRows := gethlyleaddresses.AddGethAddressToMockRows(mock, gethAddressResults)
	assetID := TestData1.AssetID

	mock.ExpectQuery("^WITH sender_table as").WithArgs(*assetID).WillReturnRows(mockRows)
	foundAddresses, err := GetDistinctAddressesFromAssetId(mock, assetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetDistinctAddressesFromAssetId", err)
	}
	for i, sourceGethAddress := range gethAddressResults {
		if cmp.Equal(sourceGethAddress, foundAddresses[i]) == false {
			t.Errorf("Expected foundTxnHashes From Method GetDistinctAddressesFromAssetId: %v is different from actual %v", sourceGethAddress, foundAddresses[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctAddressesFromAssetIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	mock.ExpectQuery("^WITH sender_table as").WithArgs(assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetDistinctAddressesFromAssetId(mock, &assetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDistinctAddressesFromAssetId", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetDistinctAddressesFromAssetId: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctTransactionHashesFromAssetIdAndStartingBlock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHashResults := []string{TestData1.TxnHash, TestData2.TxnHash}
	mockRows := mock.NewRows([]string{"txn_hash"}).AddRow(txnHashResults[0]).AddRow(txnHashResults[1])
	assetID := TestData1.AssetID
	startingBlock := TestData1.BlockNumber

	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*startingBlock, *assetID).WillReturnRows(mockRows)
	foundTxnHashes, err := GetDistinctTransactionHashesFromAssetIdAndStartingBlock(mock, assetID, startingBlock)
	if err != nil {
		t.Fatalf("an error '%s' in GetDistinctTransactionHashesFromAssetIdAndStartingBlock", err)
	}
	for i, sourceTxnHash := range txnHashResults {
		if cmp.Equal(sourceTxnHash, foundTxnHashes[i]) == false {
			t.Errorf("Expected foundTxnHashes From Method GetDistinctTransactionHashesFromAssetIdAndStartingBlock: %v is different from actual %v", sourceTxnHash, foundTxnHashes[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctTransactionHashesFromAssetIdAndStartingBlockForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	assetID := -1
	startingBlock := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*startingBlock, *&assetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetDistinctTransactionHashesFromAssetIdAndStartingBlock(mock, &assetID, startingBlock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetDistinctTransactionHashesFromAssetIdAndStartingBlock", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetDistinctTransactionHashesFromAssetIdAndStartingBlock: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

// next
func TestGetHighestBlockFromBaseAssetId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	baseAssetID := targetData.BaseAssetID
	highestBlockNumber := targetData.BlockNumber
	mockRows := mock.NewRows([]string{"block_number"}).AddRow(*highestBlockNumber)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	foundBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if cmp.Equal(*foundBlockNumber, *highestBlockNumber) == false {
		t.Errorf("Expected foundBlockNumber From Method GetHighestBlockFromBaseAssetId: %v is different from actual %v", foundBlockNumber, highestBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetHighestBlockFromBaseAssetIdForErrNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData2
	baseAssetID := targetData.BaseAssetID
	noRows := pgxmock.NewRows(DBColumns)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*baseAssetID).WillReturnRows(noRows)
	foundBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if *foundBlockNumber != uint64(0) {
		t.Errorf("Expected foundBlockNumber From Method GetHighestBlockFromBaseAssetId: to be 0 but got this: %v", foundBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetHighestBlockFromBaseAssetIdForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundBlockNumber, err := GetHighestBlockFromBaseAssetId(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetHighestBlockFromBaseAssetId", err)
	}
	if foundBlockNumber != nil {
		t.Errorf("Expected GethTransfer From Method GetHighestBlockFromBaseAssetId: to be empty but got this: %v", foundBlockNumber)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromTokenAddress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransfer{TestData1, TestData2}
	mockRows := AddGethTransferToMockRows(mock, dataList)
	tokenAddressID := TestData1.TokenAddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*tokenAddressID).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferByFromTokenAddress(mock, tokenAddressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByFromTokenAddress", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransferByFromTokenAddress: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromTokenAddressForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tokenAddressID := -999
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(tokenAddressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferByFromTokenAddress(mock, &tokenAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferByFromTokenAddress", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferByFromTokenAddress: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromMakerAddressAndTokenAddressID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransfer{TestData1, TestData2}
	mockRows := AddGethTransferToMockRows(mock, dataList)
	tokenAddressID := TestData1.TokenAddressID
	makerAddressID := TestData1.SenderAddressID
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*tokenAddressID, *makerAddressID).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferByFromMakerAddressAndTokenAddressID(mock, makerAddressID, tokenAddressID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByFromMakerAddressAndTokenAddressID", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransferByFromMakerAddressAndTokenAddressID: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromMakerAddressAndTokenAddressIDForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	tokenAddressID := -999
	makerAddressID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(tokenAddressID, makerAddressID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferByFromMakerAddressAndTokenAddressID(mock, &makerAddressID, &tokenAddressID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferByFromMakerAddressAndTokenAddressID", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferByFromMakerAddressAndTokenAddressID: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransfer{TestData1, TestData2}
	mockRows := AddGethTransferToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	makerAddressID := TestData1.SenderAddressID
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*baseAssetID, *makerAddressID, *blockNumber).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber(mock, makerAddressID, baseAssetID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumberForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -999
	makerAddressID := 1
	blockNumber := uint64(10000)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(baseAssetID, makerAddressID, blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber(mock, &makerAddressID, &baseAssetID, &blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferByFromMakerAddressAndTokenAddressIDAndBeforeBlockNumber: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromBaseAssetIDAndBeforeBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransfer{TestData1, TestData2}
	mockRows := AddGethTransferToMockRows(mock, dataList)
	baseAssetID := TestData1.BaseAssetID
	blockNumber := TestData1.BlockNumber
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(*baseAssetID, *blockNumber).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber(mock, baseAssetID, blockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferByFromBaseAssetIDAndBeforeBlockNumberForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := 1
	blockNumber := uint64(10000)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(baseAssetID, blockNumber).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber(mock, &baseAssetID, &blockNumber)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferByFromBaseAssetIDAndBeforeBlockNumber: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransfersByTxnHash(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []GethTransfer{TestData1, TestData2}
	mockRows := AddGethTransferToMockRows(mock, dataList)
	txnHash := TestData1.TxnHash
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(txnHash, *baseAssetID).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransfersByTxnHash(mock, txnHash, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransfersByTxnHash", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransfersByTxnHash: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransfersByTxnHashForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	txnHash := "test-txn-hash"
	baseAssetID := 1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(txnHash, baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransfersByTxnHash(mock, txnHash, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransfersByTxnHash", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransfersByTxnHash: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransfersByTxnHashes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransferToMockRows(mock, dataList)
	txnHashes := []string{TestData1.TxnHash, TestData2.TxnHash}
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(pq.Array(txnHashes), *baseAssetID).WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransfersByTxnHashes(mock, txnHashes, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransfersByTxnHashes", err)
	}
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransfersByTxnHash", err)
	}
	for i, sourceGethTransfer := range dataList {
		if cmp.Equal(sourceGethTransfer, foundGethTransferList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransfersByTxnHash: %v is different from actual %v", sourceGethTransfer, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransfersByTxnHashesForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()

	txnHashes := []string{"0x-invalid-1", "0x-invalid-2"}
	baseAssetID := -1
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WithArgs(pq.Array(txnHashes), baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransfersByTxnHashes(mock, txnHashes, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransfersByTxnHashes", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransfersByTxnHashes: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestRemoveGethTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	gethTransferID := targetData.ID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(*gethTransferID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransfer(mock, gethTransferID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	gethTransferID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(gethTransferID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransfer(mock, &gethTransferID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransfersFromBaseAssetIDAndStartBlockNumber(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	baseAssetID := targetData.BaseAssetID
	startBlockNumber := targetData.BlockNumber
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(*baseAssetID, *startBlockNumber).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber(mock, baseAssetID, startBlockNumber)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransfersFromBaseAssetIDAndStartBlockNumberOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	startBlockNumber := uint64(10000)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(baseAssetID, startBlockNumber).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransfersFromBaseAssetIDAndStartBlockNumber(mock, &baseAssetID, &startBlockNumber)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransfersFromBaseAssetID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	baseAssetID := targetData.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveGethTransfersFromBaseAssetID(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveGethTransfersFromBaseAssetID", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestRemoveGethTransfersFromBaseAssetIDOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM geth_transfers").WithArgs(baseAssetID).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))
	mock.ExpectRollback()
	err = RemoveGethTransfersFromBaseAssetID(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData

	mockRows := AddGethTransferToMockRows(mock, dataList)
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferList(mock)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferList", err)
	}
	testMarketDataList := TestAllData
	for i, foundGethTransfer := range foundGethTransferList {
		if cmp.Equal(foundGethTransfer, testMarketDataList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetGethTransferList: %v is different from actual %v", foundGethTransfer, testMarketDataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferListForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferList(mock)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferList", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferList: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(
		targetData.ChainID,             //1
		targetData.TokenAddress,        //2
		targetData.TokenAddressID,      //3
		targetData.AssetID,             //4
		targetData.BlockNumber,         //5
		targetData.IndexNumber,         //6
		targetData.TransferDate,        //7
		targetData.TxnHash,             //8
		targetData.SenderAddress,       //9
		targetData.SenderAddressID,     //10
		targetData.ToAddress,           //11
		targetData.ToAddressID,         //12
		targetData.Amount,              //13
		targetData.Description,         //14
		targetData.UpdatedBy,           //15
		targetData.GethProcessJobID,    //16
		pq.Array(targetData.TopicsStr), //17
		targetData.StatusID,            //18
		targetData.BaseAssetID,         //19
		targetData.TransferTypeID,      //20
		targetData.ID,                  //21
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransfer(mock, &targetData)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransferOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransfer(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	targetData.ID = utils.Ptr[int](-1)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(
		targetData.ChainID,             //1
		targetData.TokenAddress,        //2
		targetData.TokenAddressID,      //3
		targetData.AssetID,             //4
		targetData.BlockNumber,         //5
		targetData.IndexNumber,         //6
		targetData.TransferDate,        //7
		targetData.TxnHash,             //8
		targetData.SenderAddress,       //9
		targetData.SenderAddressID,     //10
		targetData.ToAddress,           //11
		targetData.ToAddressID,         //12
		targetData.Amount,              //13
		targetData.Description,         //14
		targetData.UpdatedBy,           //15
		targetData.GethProcessJobID,    //16
		pq.Array(targetData.TopicsStr), //17
		targetData.StatusID,            //18
		targetData.BaseAssetID,         //19
		targetData.TransferTypeID,      //20
		targetData.ID,                  //21
	).WillReturnError(fmt.Errorf("Cannot have -1 as ID"))

	mock.ExpectRollback()
	err = UpdateGethTransfer(mock, &targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransfer(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transfers").WithArgs(
		targetData.ChainID,             //1
		targetData.TokenAddress,        //2
		targetData.TokenAddressID,      //3
		targetData.AssetID,             //4
		targetData.BlockNumber,         //5
		targetData.IndexNumber,         //6
		targetData.TransferDate,        //7
		targetData.TxnHash,             //8
		targetData.SenderAddress,       //9
		targetData.SenderAddressID,     //10
		targetData.ToAddress,           //11
		targetData.ToAddressID,         //12
		targetData.Amount,              //13
		targetData.Description,         //14
		targetData.CreatedBy,           //15
		targetData.GethProcessJobID,    //16
		pq.Array(targetData.TopicsStr), //17
		targetData.StatusID,            //18
		targetData.BaseAssetID,         //19
		targetData.TransferTypeID,      //20
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit()
	gethTransferID, newUUID, err := InsertGethTransfer(mock, &targetData)
	if gethTransferID < 0 {
		t.Fatalf("gethTransferID should not be negative ID: %d", gethTransferID)
	}
	if newUUID == "" {
		t.Fatalf("newUUID should not be empty string: %s", newUUID)
	}
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransfer", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransferOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	targetData.BlockNumber = nil
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transfers").WithArgs(
		targetData.ChainID,             //1
		targetData.TokenAddress,        //2
		targetData.TokenAddressID,      //3
		targetData.AssetID,             //4
		targetData.BlockNumber,         //5
		targetData.IndexNumber,         //6
		targetData.TransferDate,        //7
		targetData.TxnHash,             //8
		targetData.SenderAddress,       //9
		targetData.SenderAddressID,     //10
		targetData.ToAddress,           //11
		targetData.ToAddressID,         //12
		targetData.Amount,              //13
		targetData.Description,         //14
		targetData.CreatedBy,           //15
		targetData.GethProcessJobID,    //16
		pq.Array(targetData.TopicsStr), //17
		targetData.StatusID,            //18
		targetData.BaseAssetID,         //19
		targetData.TransferTypeID,      //20
	).WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTransferID, newUUID, err := InsertGethTransfer(mock, &targetData)
	if gethTransferID >= 0 {
		t.Fatalf("Expecting -1 for ID because of error gethTransferID: %d", gethTransferID)
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

func TestInsertGethTransferOnFailureOnCommit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	targetData := TestData1
	// name can't be nil
	uuid := "01ef85e8-2c26-441e-8c7f-71d79518ad72"
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO geth_transfers").WithArgs(
		targetData.ChainID,             //1
		targetData.TokenAddress,        //2
		targetData.TokenAddressID,      //3
		targetData.AssetID,             //4
		targetData.BlockNumber,         //5
		targetData.IndexNumber,         //6
		targetData.TransferDate,        //7
		targetData.TxnHash,             //8
		targetData.SenderAddress,       //9
		targetData.SenderAddressID,     //10
		targetData.ToAddress,           //11
		targetData.ToAddressID,         //12
		targetData.Amount,              //13
		targetData.Description,         //14
		targetData.CreatedBy,           //15
		targetData.GethProcessJobID,    //16
		pq.Array(targetData.TopicsStr), //17
		targetData.StatusID,            //18
		targetData.BaseAssetID,         //19
		targetData.TransferTypeID,      //20
	).WillReturnRows(pgxmock.NewRows([]string{"id", "uuid"}).AddRow(1, uuid))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("Random SQL Error"))
	mock.ExpectRollback()
	gethTransferID, newUUID, err := InsertGethTransfer(mock, &targetData)
	if gethTransferID >= 0 {
		t.Fatalf("Expecting -1 for gethTransferID because of error gethTransferID: %d", gethTransferID)
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

func TestInsertGethTransfers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transfers"}, DBColumnsInsertGethTransfers)
	targetData := TestAllData
	err = InsertGethTransfers(mock, targetData)
	if err != nil {
		t.Fatalf("an error '%s' in InsertGethTransfers", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestInsertGethTransfersOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"geth_transfers"}, DBColumnsInsertGethTransfers).WillReturnError(fmt.Errorf("Random SQL Error"))
	targetData := TestAllData
	err = InsertGethTransfers(mock, targetData)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransferAddresses(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	baseAssetID := TestData1.BaseAssetID
	defer mock.Close()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransferAddresses(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransferAddresses", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransferAddressesOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransferAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransferAddressesOnFirstFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnError(fmt.Errorf("1st SQL Error"))

	mock.ExpectRollback()
	err = UpdateGethTransferAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransferAddressesOnSecondFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnError(fmt.Errorf("2nd SQL Error"))

	mock.ExpectRollback()
	err = UpdateGethTransferAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransferAddressesOnThirdFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("^UPDATE geth_transfers").WithArgs(*baseAssetID).WillReturnError(fmt.Errorf("3rd SQL Error"))

	mock.ExpectRollback()
	err = UpdateGethTransferAddresses(mock, baseAssetID)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromTransfers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := []string{TestData1.SenderAddress, TestData2.ToAddress}
	mockRows := mock.NewRows([]string{"address"}).AddRow(dataList[0]).AddRow(dataList[1])
	baseAssetID := TestData1.BaseAssetID
	mock.ExpectQuery("^WITH sender_table as ").WithArgs(*baseAssetID).WillReturnRows(mockRows)
	foundNullAddresses, err := GetNullAddressStrsFromTransfers(mock, baseAssetID)
	if err != nil {
		t.Fatalf("an error '%s' in GetNullAddressStrsFromTransfers", err)
	}
	for i, nullAddress := range foundNullAddresses {
		if cmp.Equal(nullAddress, dataList[i]) == false {
			t.Errorf("Expected GethTransfer From Method GetNullAddressStrsFromTransfers: %v is different from actual %v", nullAddress, dataList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetNullAddressStrsFromTransfersForErr(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	baseAssetID := -1
	mock.ExpectQuery("^WITH sender_table as ").WithArgs(baseAssetID).WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetNullAddressStrsFromTransfers(mock, &baseAssetID)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetNullAddressStrsFromTransfers", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetNullAddressStrsFromTransfers: to be empty but got this: %v", foundGethTransferList)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestUpdateGethTransfersAssetIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateGethTransfersAssetIDs(mock)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateGethTransfersAssetIDs", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransfersAssetIDsOnFailureAtBegin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("Failure at begin"))
	err = UpdateGethTransfersAssetIDs(mock)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}
func TestUpdateGethTransfersAssetIDsOnFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE geth_transfers").WillReturnError(fmt.Errorf("1st SQL Error"))

	mock.ExpectRollback()
	err = UpdateGethTransfersAssetIDs(mock)
	if err == nil {
		t.Fatalf("was expecting an error, but there was none")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer mock.Close()
	dataList := TestAllData
	mockRows := AddGethTransferToMockRows(mock, dataList)
	_start := 0
	_end := 10
	_sort := "id"
	_order := "ASC"
	filters := []string{"chain_id = 1"}
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WillReturnRows(mockRows)
	foundGethTransferList, err := GetGethTransferListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetGethTransferListByPagination", err)
	}
	for i, sourceData := range dataList {
		if cmp.Equal(sourceData, foundGethTransferList[i]) == false {
			t.Errorf("Expected sourceData From Method GetGethTransferListByPagination: %v is different from actual %v", sourceData, foundGethTransferList[i])
		}
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There awere unfulfilled expectations: %s", err)
	}
}

func TestGetGethTransferListByPaginationForErr(t *testing.T) {
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
	mock.ExpectQuery("^SELECT (.+) FROM geth_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
	foundGethTransferList, err := GetGethTransferListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err == nil {
		t.Fatalf("expected an error '%s' in GetGethTransferListByPagination", err)
	}
	if len(foundGethTransferList) != 0 {
		t.Errorf("Expected From Method GetGethTransferListByPagination: to be empty but got this: %v", foundGethTransferList)
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transfers").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(numOfChainsExpected))
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
	mock.ExpectQuery("^SELECT COUNT(.*) FROM geth_transfers").WillReturnError(pgx.ScanArgError{Err: errors.New("Random SQL Error")})
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
