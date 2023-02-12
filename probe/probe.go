package probe

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	queryTemplate = `SELECT "%s" FROM "%s" WHERE "%s" = "%s"`
	queryTemplateExt = `SELECT "%s" FROM "%s" WHERE "%s" = "%s" AND "%s" = "%s"`
)

var tempQueries = map[string][]string{
	"first_name": 		  {"DataValue", "MiscDataDynamic", "DataName", "PlayerFirstName"},
	"surname":    		  {"DataValue", "MiscDataDynamic", "DataName", "PlayerLastName"},
	"xp":         		  {"DataValue", "MiscDataDynamic", "DataName", "ExperiencePoints"},
	"house":      		  {"DataValue", "MiscDataDynamic", "DataName", "HouseID"},
	"galleons":   		  {"Count", "InventoryDynamic", "CharacterID", "Player0", "ItemID", "Knuts"},
	"inventory_size":     {"DataValue", "MiscDataDynamic", "DataName", "BaseInventoryCapacity"},
	"talent_points":      {"DataValue", "MiscDataDynamic", "DataOwner", "Player0", "DataName", "PerkPoints"},
}

var queries = map[string]string{
	"inventory": `SELECT "ItemID", "Count" FROM "InventoryDynamic" WHERE "CharacterID" = "Player0" AND "HolderID" = "ResourceInventory" AND "ItemID" IS NOT NULL`,
}

func queryRow(db *sql.DB, q string) (*sql.Row, error) {
	row := db.QueryRow(q)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	return row, nil
}

func readStr(db *sql.DB, q string) (string, error) {
	var val string
	row, err := queryRow(db, q)
	if err != nil {
		return "", err
	}
	err = row.Scan(&val)
	return val, err
}


func readInt(db *sql.DB, q string) (int64, error) {
	var val int64
	row, err := queryRow(db, q)
	if err != nil {
		return 0, err
	}
	err = row.Scan(&val)
	return val, err
}

func makeQuery(key string) string {
	queryStrings := tempQueries[key]

	var q string
	if len(queryStrings) > 4 {
		q = fmt.Sprintf(
			queryTemplateExt, queryStrings[0], queryStrings[1],
			queryStrings[2], queryStrings[3], queryStrings[4],
			queryStrings[5],
		)
	} else {
		q = fmt.Sprintf(
			queryTemplate, queryStrings[0], queryStrings[1],
			queryStrings[2], queryStrings[3],
		)
	}
	return q
}

func readPlayerName(db *sql.DB) (string, error) {
	var playerName string
	for idx, queryStr := range [2]string{"first_name", "surname"} {
		q := makeQuery(queryStr)
		val, err := readStr(db, q)
		if err != nil {
			return "", err
		}
		playerName += val
		if idx == 0 {
			playerName += " "
		}
	}

	return playerName, nil
}

func Run(db *sql.DB) error {
	playerName, err := readPlayerName(db)
	if err != nil {
		return err
	}

	xp, err := readInt(db, makeQuery("xp"))
	if err != nil {
		return err
	}

	house, err := readStr(db, makeQuery("house"))
	if err != nil {
		return err
	}

	galleons, err := readInt(db, makeQuery("galleons"))
	if err != nil {
		return err
	}

	inv_size, err := readInt(db, makeQuery("inventory_size"))
	if err != nil {
		return err
	}

	talentPoints, err := readInt(db, makeQuery("talent_points"))
	if err != nil {
		if err == sql.ErrNoRows {
			talentPoints = 0
		} else {
			return err
		}
	}

	fmt.Printf(
		"Name:           %s\nXP:             %d\nHouse:          %s\nGalleons:       %d\nInventory size: %d\nTalent points:  %d\n\n",
		playerName, xp, house, galleons, inv_size, talentPoints,
	)

	fmt.Println("Inventory:")
	var itemId string
	var quantity int64

	rows, err := db.Query(queries["inventory"])
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemId, &quantity)
		if err != nil {
			return err
		}
		fmt.Printf("%s, %d\n", itemId, quantity)
	}
	return nil
}