package assetchainlinkfeed

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var TestFeed1 = AssetChainLinkFeed{
	AssetID:                          ptr(1),
	ChainID:                          ptr(1),
	ChainlinkDataFeedContractAddress: "0x1234567890abcdef",
	CreatedBy:                        "SYSTEM",
	CreatedAt:                        time.Now(),
	UpdatedBy:                        "SYSTEM",
	UpdatedAt:                        time.Now(),
}

var TestFeed2 = AssetChainLinkFeed{
	AssetID:                          ptr(2),
	ChainID:                          ptr(1),
	ChainlinkDataFeedContractAddress: "0xabcdef1234567890",
	CreatedBy:                        "SYSTEM",
	CreatedAt:                        time.Now(),
	UpdatedBy:                        "SYSTEM",
	UpdatedAt:                        time.Now(),
}

var TestFeeds = []AssetChainLinkFeed{TestFeed1, TestFeed2}

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

func TestGetAssetChainLinkFeed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mockRows := AddAssetChainLinkFeedToMockRows(mock, []AssetChainLinkFeed{target})
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WithArgs(*target.AssetID, *target.ChainID).WillReturnRows(mockRows)
	found, err := GetAssetChainLinkFeed(mock, target.AssetID, target.ChainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainLinkFeed", err)
	}
	if found == nil || *found.AssetID != *target.AssetID || *found.ChainID != *target.ChainID {
		t.Errorf("Expected AssetChainLinkFeed: %v, got: %v", target, found)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainLinkFeedForNoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	assetID := 999
	chainID := 999
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WithArgs(assetID, chainID).WillReturnRows(pgxmock.NewRows(columns))
	found, err := GetAssetChainLinkFeed(mock, &assetID, &chainID)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainLinkFeed", err)
	}
	if found != nil {
		t.Errorf("Expected nil, got: %v", found)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainLinkFeedList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetChainLinkFeedToMockRows(mock, TestFeeds)
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WillReturnRows(mockRows)
	feeds, err := GetAssetChainLinkFeedList(mock, []int{1, 2}, []int{1})
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainLinkFeedList", err)
	}
	if len(feeds) != len(TestFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(TestFeeds), len(feeds))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestInsertAssetChainLinkFeed(t *testing.T) {
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
	err = InsertAssetChainLinkFeed(mock, &target)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetChainLinkFeed", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateAssetChainLinkFeed(t *testing.T) {
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
	err = UpdateAssetChainLinkFeed(mock, &target)
	if err != nil {
		t.Fatalf("an error '%s' in UpdateAssetChainLinkFeed", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRemoveAssetChainLinkFeed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	target := TestFeed1
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM asset_chain_link_data_feed").WithArgs(*target.AssetID, *target.ChainID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err = RemoveAssetChainLinkFeed(mock, target.AssetID, target.ChainID)
	if err != nil {
		t.Fatalf("an error '%s' in RemoveAssetChainLinkFeed", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func AddAssetChainLinkFeedToMockRows(mock pgxmock.PgxPoolIface, dataList []AssetChainLinkFeed) *pgxmock.Rows {
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

func TestInsertAssetChainLinkFeeds(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mock.ExpectCopyFrom(pgx.Identifier{"asset_chain_link_data_feed"}, columns)
	err = InsertAssetChainLinkFeeds(mock, TestFeeds)
	if err != nil {
		t.Fatalf("an error '%s' in InsertAssetChainLinkFeeds", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetChainLinkFeedListByPagination(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	mockRows := AddAssetChainLinkFeedToMockRows(mock, TestFeeds)
	mock.ExpectQuery("^SELECT (.+) FROM asset_chain_link_data_feed").WillReturnRows(mockRows)
	_start := 0
	_end := 10
	_order := "asset_id"
	_sort := "ASC"
	filters := []string{}
	feeds, err := GetAssetChainLinkFeedListByPagination(mock, &_start, &_end, _order, _sort, filters)
	if err != nil {
		t.Fatalf("an error '%s' in GetAssetChainLinkFeedListByPagination", err)
	}
	if len(feeds) != len(TestFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(TestFeeds), len(feeds))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
