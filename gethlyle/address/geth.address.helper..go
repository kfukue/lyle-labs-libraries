package gethlyleaddresses

import (
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
	"github.com/pashagolub/pgxmock/v4"
)

var DBColumns = []string{
	"id",              //1
	"uuid",            //2
	"name",            //3
	"alternate_name",  //4
	"description",     //5
	"address_str",     //6
	"address_type_id", //7
	"created_by",      //8
	"created_at",      //9
	"updated_by",      //10
	"updated_at",      //11
}
var DBColumnsInsertGethAddressList = []string{
	"uuid",            //1
	"name",            //2
	"alternate_name",  //3
	"description",     //4
	"address_str",     //5
	"address_type_id", //6
	"created_by",      //7
	"created_at",      //8
	"updated_by",      //9
	"updated_at",      //10
}

var TestData1 = GethAddress{
	ID:            utils.Ptr[int](1),
	UUID:          "880607ab-2833-4ad7-a231-b983a61c7b39",
	Name:          "EOA: 0xFC3d170c29581E60861Ac2b500b098722d9861e9",
	AlternateName: "EOA: 0xFC3d170c29581E60861Ac2b500b098722d9861e9",
	Description:   "",
	AddressStr:    "0xFC3d170c29581E60861Ac2b500b098722d9861e9",
	AddressTypeID: utils.Ptr[int](utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID),
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}

var TestData2 = GethAddress{
	ID:            utils.Ptr[int](2),
	UUID:          "880607ab-2833-4ad7-a231-b983a61cad34",
	Name:          "Contract: 0x40762e9b87aa6457f069925a86352d13339cb68f",
	AlternateName: "Contract: 0x40762e9b87aa6457f069925a86352d13339cb68f",
	Description:   "",
	AddressStr:    "0x40762e9b87aa6457f069925a86352d13339cb68f",
	AddressTypeID: utils.Ptr[int](utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_ID),
	CreatedBy:     "SYSTEM",
	CreatedAt:     utils.SampleCreatedAtTime,
	UpdatedBy:     "SYSTEM",
	UpdatedAt:     utils.SampleCreatedAtTime,
}
var TestAllData = []GethAddress{TestData1, TestData2}

func AddGethAddressToMockRows(mock pgxmock.PgxPoolIface, dataList []GethAddress) *pgxmock.Rows {
	rows := mock.NewRows(DBColumns)
	for _, data := range dataList {
		rows.AddRow(
			data.ID,            //1
			data.UUID,          //2
			data.Name,          //3
			data.AlternateName, //4
			data.Description,   //5
			data.AddressStr,    //6
			data.AddressTypeID, //7
			data.CreatedBy,     //8
			data.CreatedAt,     //9
			data.UpdatedBy,     //10
			data.UpdatedAt,     //11
		)
	}
	return rows
}
