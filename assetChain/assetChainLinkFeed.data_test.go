package assetchain

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var TestFeed1 = AssetChain{
	AssetID:                          ptr(1),
	ChainID:                          ptr(1),
	ChainlinkDataFeedContractAddress: "0x1234567890abcdef",
	CreatedBy:                        "SYSTEM",
	CreatedAt:                        time.Now(),
	UpdatedBy:                        "SYSTEM",
	UpdatedAt:                        time.Now(),
}

var TestFeed2 = AssetChain{
	AssetID:                          ptr(2),
	ChainID:                          ptr(1),
	ChainlinkDataFeedContractAddress: "0xabcdef1234567890",
	CreatedBy:                        "SYSTEM",
	CreatedAt:                        time.Now(),
	UpdatedBy:                        "SYSTEM",
	UpdatedAt:                        time.Now(),
}

var TestFeeds = []AssetChain{TestFeed1, TestFeed2}

var columns = []string{
	"asset_id",
	"chain_id",
	"chainlink_data_feed_contract_address",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
}

func ptr(i int) *int { return &i }

func TestGetAssetChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mockRows := AddAssetChainToMockRows(mock, []AssetChain{target})
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WithArgs(*target.AssetID, *target.ChainID).WillReturnRows(mockRows)
	found, err := GetAssetChain(mock, target.AssetID, target.ChainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChain", err)
	}
	if found == nil || *found.AssetID != *target.AssetID || *found.ChainID != *target.ChainID {
		t.Errorf("Expected AssetChain: %v, got: %v", target, found)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainForNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	assetID := 999
	chainID := 999
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WithArgs(assetID, chainID).WillReturnRows(pgxmock.NewRows(columns))
	found, err := GetAssetChain(mock, &assetID, &chainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChain", err)
	}
	if found != nil {
		t.Errorf("Expected nil, got: %v", found)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetChainToMockRows(mock, TestFeeds)
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WillReturnRows(mockRows)
	feeds, err := GetAssetChainList(mock, []int{1, 2}, []int{1})
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainList", err)
	}
	if len(feeds) != len(TestFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(TestFeeds), len(feeds))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO asset_chain_link_data_feed").WithArgs(
		target.AssetID,
		target.ChainID,
		target.ChainlinkDataFeedContractAddress,
		target.CreatedBy,
		target.UpdatedBy,
	).WillReturnRows(pgxmock.NewRows([]string{"asset_id", "chain_id"}).AddRow(*target.AssetID, *target.ChainID))
	mock.ExpectCommit()
	err = InsertAssetChain(mock, &target)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE asset_chain_link_data_feed").WithArgs(
		target.ChainlinkDataFeedContractAddress,
		target.UpdatedBy,
		target.AssetID,
		target.ChainID,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err = UpdateAssetChain(mock, &target)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAssetChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetChain(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_chain_link_data_feed").WithArgs(*target.AssetID, *target.ChainID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAssetChain(mock, target.AssetID, target.ChainID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAssetChain", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func AddAssetChainToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetChain) *pgxmock.Rows {
	rows := mock.NewRows(columns)
	for _, data := range dataList {
		rows.AddRow(
			data.AssetID,
			data.ChainID,
			data.ChainlinkDataFeedContractAddress,
			data.CreatedBy,
			data.CreatedAt,
			data.UpdatedBy,
			data.UpdatedAt,
		)
	}
	return rows
}

func TestInsertAssetChains(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"asset_chain_link_data_feed"}, columns)
	err = InsertAssetChains(mock, TestFeeds)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetChains", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetChainToMockRows(mock, TestFeeds)
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WillReturnRows(mockRows)
	_start := 0
	_end := 10
	_order := "asset_id"
	_sort := "ASC"
	filters := []string{}
	feeds, err := GetAssetChainListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainListByPagination", err)
	}
	if len(feeds) != len(TestFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(TestFeeds), len(feeds))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
