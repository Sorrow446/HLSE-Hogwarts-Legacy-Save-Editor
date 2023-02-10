package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"database/sql"
	"path/filepath"
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/alexflint/go-arg"
	_ "github.com/mattn/go-sqlite3"

	"main/probe"
)

var dbPath = filepath.Join(os.TempDir(), "hlse_tmp.db")

var magic = [4]byte{'\x47', '\x56', '\x41', '\x53'}

var rawDbImageStr = []byte{
	'\x52', '\x61', '\x77', '\x44', '\x61', '\x74',
	'\x61', '\x62', '\x61', '\x73', '\x65', '\x49',
	'\x6D', '\x61', '\x67', '\x65',
}

func extractDb(saveData []byte) (int, int, error) {
	imageStrStart := bytes.Index(saveData, rawDbImageStr)
	if imageStrStart == -1 {
		return 0, 0, errors.New("couldn't find db image string")
	}
	dbSizeOffset := imageStrStart+61
	dbStartOffset := dbSizeOffset+4
	dbSizeBytes := saveData[dbSizeOffset:dbStartOffset]
	dbSize := binary.LittleEndian.Uint32(dbSizeBytes)
	dbEndOffset := dbStartOffset+int(dbSize)
	dbData := saveData[dbStartOffset:dbEndOffset]
	err := os.WriteFile(dbPath, dbData, 0755)
	return imageStrStart, dbEndOffset, err
}

var queries = map[string]string{
	"galleons":        `UPDATE "InventoryDynamic" SET "Count" = %d WHERE "CharacterID" = "Player0" AND "ItemID" = "Knuts"`,
	"xp":              `UPDATE "MiscDataDynamic" SET "DataValue" = %d WHERE "DataName" = "ExperiencePoints"`,
	"first_name":      `UPDATE "MiscDataDynamic" SET "DataValue" = "%s" WHERE "DataName" = "PlayerFirstName"`,
	"surname":         `UPDATE "MiscDataDynamic" SET "DataValue" = "%s" WHERE "DataName" = "PlayerLastName"`,
	"inventory_size":  `UPDATE "MiscDataDynamic" SET "DataValue" = %d WHERE "DataName" = "BaseInventoryCapacity"`,
	"inventory_quant": `UPDATE "InventoryDynamic" SET "Count" = %d WHERE "CharacterID" = "Player0" AND "HolderID" = "ResourceInventory" AND LOWER(ItemId) = "%s"`,
}

func containsItemId(parsed []*ItemPairs, itemId string) bool {
	for _, pair := range parsed {
		if pair.ItemID == itemId {
			return true
		}
	}
	return false
}

func parseInvPairs(pairs []string) ([]*ItemPairs, error) {
	pairsLen := len(pairs)
	if pairsLen%2 !=0 {
		return nil, errors.New("item quantity pairs can't be odd")
	}

	var parsed []*ItemPairs

	for i := 0; i < pairsLen; i+=2 {
		itemId := strings.ToLower(pairs[i])
		quantity, err := strconv.ParseInt(pairs[i+1], 10, 64)
		if err != nil {
		    return nil, err
		}
		pair := &ItemPairs{
			ItemID: itemId,
			Quantity: quantity,
		}
		if containsItemId(parsed, itemId) {
			fmt.Println("filtered pair with same item ID:", pair)
			continue
		}

		parsed = append(parsed, pair)
	}

	return parsed, nil
}

func parseArgs() (*Args, error) {
	var args Args
	arg.MustParse(&args)
	if len(args.ItemQuantities) > 0 {
		invPairs, err := parseInvPairs(args.ItemQuantities)
		if err != nil {
			return nil, err
		}
		args.ParsedItemQuants = invPairs
	}
	if args.XP < 0 {
		return nil, errors.New("xp can't be negative")
	}
	if args.Galleons < 0 {
		return nil, errors.New("galleons can't be negative")
	}

	if args.InventorySize != 0 && args.InventorySize < 20 {
		return nil, errors.New("inventory size can't be less than 20")
	}

	if !args.Probe && args.XP == 0 && args.Galleons == 0 && args.InventorySize == 0 && len(args.ItemQuantities) < 1 && args.FirstName == "" && args.Surname == "" {
		return nil, errors.New("no write arguments were provided")
	}

	if args.OutPath == "" {
		args.OutPath = args.InPath
	}

	return &args, nil
}

func updateRow(db *sql.DB, q string) error {
	res, err := db.Exec(q)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
	 	return err
	}
	if rowsAffected == 0 {
		return errors.New("db row wasn't updated")
	}
	return nil
}


func writeSave(updatedDbBytes, saveData []byte, imageStrStart, dbEndOffset int, outPath string) error {
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(saveData[:imageStrStart+35])
	if err != nil {
		return err
	}
	buf := make([]byte, 4)


	updatedDbSize := len(updatedDbBytes)
	binary.LittleEndian.PutUint32(buf, uint32(updatedDbSize+4))

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	_, err = f.Write(saveData[imageStrStart+39:imageStrStart+61])
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(buf, uint32(updatedDbSize))

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	_, err = f.Write(updatedDbBytes)
	if err != nil {
		return err
	}
	_, err = f.Write(saveData[dbEndOffset:])
	return err
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	saveData, err := os.ReadFile(args.InPath)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(saveData[:4], magic[:]) {
		panic("invalid save file magic")
	}

	imageStrStart, dbEndOffset, err := extractDb(saveData)
	if err != nil {
		panic(err)
	}

	defer os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	if args.Probe {
		err = probe.Run(db)
		db.Close()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}
	
	if args.XP > 0 {
		err = updateRow(db, fmt.Sprintf(queries["xp"], args.XP))
		if err != nil {
			db.Close()
			panic(err)
		}
	}

	if args.Galleons > 0 {
		err = updateRow(db, fmt.Sprintf(queries["galleons"], args.Galleons))
		if err != nil {
			db.Close()
			panic(err)
		}
	}

	if args.InventorySize > 0 {
		err = updateRow(db, fmt.Sprintf(queries["inventory_size"], args.InventorySize))
		if err != nil {
			db.Close()
			panic(err)
		}
	}

	if args.FirstName != "" {
		err = updateRow(db, fmt.Sprintf(queries["first_name"], args.FirstName))
		if err != nil {
			db.Close()
			panic(err)
		}
	}

	if args.Surname != "" {
		err = updateRow(db, fmt.Sprintf(queries["surname"], args.Surname))
		if err != nil {
			db.Close()
			panic(err)
		}
	}

	if len(args.ItemQuantities) > 0 {
		for _, pair := range args.ParsedItemQuants {
			err = updateRow(db, fmt.Sprintf(queries["inventory_quant"], pair.Quantity, pair.ItemID))
			if err != nil {
				db.Close()
				panic(err)
			}
		}
	}

	db.Close()
	updatedDbBytes, err := os.ReadFile(dbPath)
	if err != nil {
		panic(err)
	}

	err = writeSave(updatedDbBytes, saveData, imageStrStart, dbEndOffset, args.OutPath)
	if err != nil {
		panic(err)
	}
}