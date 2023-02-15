package probe

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	queryTemplate = `SELECT "%s" FROM "%s" WHERE "%s" = "%s"`
	queryTemplateExt = `SELECT "%s" FROM "%s" WHERE "%s" = "%s" AND "%s" = "%s"`
)

var wandDataTypes = [5]string{"WandBase", "WandWood", "WandCore", "WandLength", "WandFlex"}

var xps = []int{
	500, 1030, 1595, 2195, 2835,
	3515, 4240, 5015, 5840, 6715,
	7650, 8650, 9700, 10825, 12025,
	13300, 14660, 16110, 17650, 19290,
	21035, 22885, 24865, 26965, 29205,
	31590, 34130, 36830, 39710, 42750,
	46000, 49500, 53000, 56500, 60000,
	63500, 67000, 70500, 74000,
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


func readI64(db *sql.DB, q string) (int64, error) {
	var val int64
	row, err := queryRow(db, q)
	if err != nil {
		return 0, err
	}
	err = row.Scan(&val)
	return val, err
}

func readInt(db *sql.DB, q string) (int, error) {
	var val int
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

func getLevelFromXp(playerXp int) int {
	for level, xp := range xps {
		if playerXp < xp {
			return level+1
		}
	}
	return 40
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

func readWand(db *sql.DB) error {
	fmt.Println("Wand:")

	for _, dataName := range wandDataTypes {
		val, err := readStr(db, fmt.Sprintf(queries["wand"], dataName))
		if err != nil {
			if err != sql.ErrNoRows && dataName != "WandStyle" {
				return err
			}
			val = "none"
		}

		fmtVal, ok := fmtWandParts[dataName][val]
		if !ok {
			fmtVal = val
		}

		fmt.Printf("%s:%s%s\n", dataName[4:], strings.Repeat(" ", 11-len(dataName)), fmtVal)
	}
	fmt.Println("")
	return nil
}

func readInv(db *sql.DB) error {
	fmt.Println("Resource inventory:")
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

func Run(db *sql.DB) error {
	playerName, err := readPlayerName(db)
	if err != nil {
		return err
	}

	xp, err := readInt(db, makeQuery("xp"))
	if err != nil {
		return err
	}
	level := getLevelFromXp(xp)

	house, err := readStr(db, makeQuery("house"))
	if err != nil {
		return err
	}

	galleons, err := readI64(db, makeQuery("galleons"))
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
		"Name:           %s\nXP:             %d\nLevel:          %d\nHouse:          %s\nGalleons:       %d\nInventory size: %d\nTalent points:  %d\n\n",
		playerName, xp, level, house, galleons, inv_size, talentPoints,
	)

	err = readWand(db)
	if err != nil {
		panic(err)
	}
	err = readInv(db)
	if err != nil {
		panic(err)
	}
	return nil
}