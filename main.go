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

var tmpDbPath = filepath.Join(os.TempDir(), "hlse_tmp.db")

var magic = []byte{'\x47', '\x56', '\x41', '\x53'}

var dbMagic = []byte{
	'\x53', '\x51', '\x4C', '\x69', '\x74', '\x65',
	'\x20', '\x66', '\x6F', '\x72', '\x6D', '\x61',
	'\x74', '\x20', '\x33',
}

var rawDbImageStr = []byte{
	'\x52', '\x61', '\x77', '\x44', '\x61', '\x74',
	'\x61', '\x62', '\x61', '\x73', '\x65', '\x49',
	'\x6D', '\x61', '\x67', '\x65',
}

var xps = []int{
	0, 500, 1030, 1595, 2195,
	2835, 3515, 4240, 5015, 5840,
	6715, 7650, 8650, 9700, 10825,
	12025, 13300, 14660, 16110, 17650, 
	19290, 21035, 22885, 24865, 26965,
	29205, 31590, 34130, 36830, 39710,
	42750, 46000, 49500, 53000, 56500,
	60000, 63500, 67000, 70500, 74000,
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


func parseWand(args *Args) (*Wand, error) {
	var wand Wand
	errTemp := "invalid wand %s"
	
	if args.WandBase != "" {
		base, ok := resolveWandParts["bases"][strings.ToLower(args.WandBase)]
		if !ok {
			return nil, fmt.Errorf(errTemp, "base")
		}
		wand.Base = base
	}

	if args.WandWood != "" {
		wood, ok := resolveWandParts["wood"][strings.ToLower(args.WandWood)]
		if !ok {
			return nil, fmt.Errorf(errTemp, "wood")
		}
		wand.Wood = wood
	}

	if args.WandCore != "" {
		core, ok := resolveWandParts["cores"][strings.ToLower(args.WandCore)]
		if !ok {
			return nil, fmt.Errorf(errTemp, "core")
		}
		wand.Core = core
	}

	if args.WandLength != "" {
		length, ok := resolveWandParts["lengths"][strings.ToLower(args.WandLength)]
		if !ok {
			return nil, fmt.Errorf(errTemp, "length")
		}
		wand.Length = length
	}

	if args.WandFlex != "" {
		flex, ok := resolveWandParts["flex"][strings.ToLower(args.WandFlex)]
		if !ok {
			return nil, fmt.Errorf(errTemp, "flex")
		}
		wand.Flex = flex
	}

	// if args.WandStyle != "" {
	// 	invalidStyle := fmt.Errorf(errTemp, "style")
	// 	fullStyle := strings.ToLower(args.WandStyle)
	// 	split := strings.SplitN(fullStyle, "_", 4)
	// 	splitLen := len(split)
	// 	if splitLen != 2  && splitLen != 4 {
	// 		return nil, invalidStyle
	// 	}
	// 	base := strings.Join(split[:2], "_")
	// 	baseRes, ok := resolveWandParts["bases"][base]
	// 	if !ok {
	// 		return nil, invalidStyle
	// 	}
	// 	if splitLen == 4 {
	// 		style := strings.Join(split[2:], "_")
	// 		styleRes, ok := resolveWandParts["styles"][style]
	// 		if !ok {
	// 			return nil, invalidStyle
	// 		}
	// 		wand.Style = baseRes + "_" + styleRes
	// 	} else {
	// 		wand.Style = baseRes
	// 	}
	// }	
	return &wand, nil
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
	if args.XP > 99999999 {
		return nil, errors.New("xp can't be more than 99999999")
	}
	if args.Galleons < 0 {
		return nil, errors.New("galleons can't be negative")
	}

	if args.InventorySize != 0 && args.InventorySize < 20 {
		return nil, errors.New("inventory size can't be less than 20")
	}

	if args.XP > 0 && args.Level > 0 {
		return nil, errors.New("xp and level args can't both be provided")
	}

	if args.Level > 0 {
		if !(args.Level >= 1 && args.Level <= 40) {
			return nil, errors.New("level must be between 1 and 40")
		} else {
			args.XP = xps[args.Level-1]
		}
	}

	noArgs := !args.Probe && !args.Unstuck && !args.DumpDB &&
		!args.InjectDB && !args.ResetTalentPoints && args.XP == 0 &&
		args.Galleons == 0 && args.InventorySize == 0 &&
		args.TalentPoints == 0 && len(args.ItemQuantities) < 1 &&
		args.Level == 0 && args.House == "" && args.FirstName == "" &&
		args.Surname == ""

	noWandArgs := args.WandBase == "" && args.WandWood  == "" &&
		args.WandCore == "" && args.WandLength == "" &&
		args.WandFlex == ""

	if noArgs && noWandArgs {
		return nil, errors.New("no write arguments were provided")
	}

	if args.TalentPoints < 0 {
		return nil, errors.New("talent points can't be negative")
	}

	if args.House != "" {
		house, ok := resolveHouse[strings.ToLower(args.House)]
		if !ok {
			return nil, errors.New("invalid house")
		}
		args.House = house
	}

	if args.DumpDB && args.OutPath == "" {
		return nil, errors.New("output path of db required when dumping db")
	}
	if args.InjectDB {
		if args.OutPath == "" {
			return nil, errors.New("output path of save file required when injecting db")
		}
		inPath := args.InPath
		outPath := args.OutPath
		args.InPath = outPath
		args.OutPath = inPath
	}

	wand, err := parseWand(&args)
	if err != nil {
		return nil, err
	}
	args.Wand = wand

	if args.OutPath == "" {
		args.OutPath = args.InPath
	}

	return &args, nil
}

func extractDb(saveData []byte, dumpDb, injectDb bool, outPath string) (string, int, int, error) {
	var (
		dbPath string
		dbData []byte
		err error
	)

	imageStrStart := bytes.Index(saveData, rawDbImageStr)
	if imageStrStart == -1 {
		return "", 0, 0, errors.New("couldn't find db image string")
	}
	dbSizeOffset := imageStrStart+61
	dbStartOffset := dbSizeOffset+4
	dbSizeBytes := saveData[dbSizeOffset:dbStartOffset]
	dbSize := binary.LittleEndian.Uint32(dbSizeBytes)
	dbEndOffset := dbStartOffset+int(dbSize)
	if injectDb {
		dbData, err = os.ReadFile(outPath)
		if err != nil {
			return "", 0, 0, err
		}
	} else {
		dbData = saveData[dbStartOffset:dbEndOffset]
	}

	if dumpDb || injectDb {
		dbPath = outPath
	} else {
		dbPath = tmpDbPath
	}

	if !injectDb {
		err = os.WriteFile(dbPath, dbData, 0755)
		if err != nil {
			return "", 0, 0, err
		}
	}
	
	return dbPath, imageStrStart, dbEndOffset, nil
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

func writeToDb(db *sql.DB, args *Args) error {
	defer db.Close()

	if args.XP > 0 || args.Level > 0 {
		err := updateRow(db, fmt.Sprintf(queries["xp"], args.XP))
		if err != nil {
			return err
		}
	}

	if args.Galleons > 0 {
		err := updateRow(db, fmt.Sprintf(queries["galleons"], args.Galleons))
		if err != nil {
			return err
		}
	}

	if args.InventorySize > 0 {
		err := updateRow(db, fmt.Sprintf(queries["inventory_size"], args.InventorySize))
		if err != nil {
			return err
		}
	}

	if args.FirstName != "" {
		err := updateRow(db, fmt.Sprintf(queries["first_name"], args.FirstName))
		if err != nil {
			return err
		}
	}

	if args.Surname != "" {
		err := updateRow(db, fmt.Sprintf(queries["surname"], args.Surname))
		if err != nil {
			return err
		}
	}

	if args.House != "" {
		err := updateRow(db, fmt.Sprintf(queries["house"], args.House))
		if err != nil {
			return err
		}
	}

	if len(args.ItemQuantities) > 0 {
		for _, pair := range args.ParsedItemQuants {
			err := updateRow(db, fmt.Sprintf(queries["inventory_quant"], pair.Quantity, pair.ItemID))
			if err != nil {
				return err
			}
		}
	}

	if args.TalentPoints > 0 || args.ResetTalentPoints {
		err := updateRow(db, fmt.Sprintf(queries["talent_points"], args.TalentPoints))
		if err != nil {
			return err
		}
	}

	if args.Unstuck {
		for dataName, dataValue := range unstuckMap {
			err := updateRow(db, fmt.Sprintf(queries["unstuck"], dataValue, dataName))
			if err != nil {
				return err
			}		
		}
		err := updateRow(db, queries["world"])
		if err != nil {
			return err
		}
	}

	wand := args.Wand

	if wand.Base != "" {
		for _, dataName := range [2]string{"WandBase", "WandStyle"} {
			err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Base, dataName))
			if err != nil {
				return err
			}
		}
	}

	if wand.Wood != "" {
		err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Wood, "WandWood"))
		if err != nil {
			return err
		}
	}

	if wand.Core != "" {
		err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Core, "WandCore"))
		if err != nil {
			return err
		}
	}

	if wand.Length != "" {
		err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Length, "WandLength"))
		if err != nil {
			return err
		}
	}

	if wand.Flex != "" {
		err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Flex, "WandFlex"))
		if err != nil {
			return err
		}
	}

	// if wand.Style != "" {
	// 	err := updateRow(db, fmt.Sprintf(queries["wand"], wand.Style, "WandStyle"))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
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

	if !bytes.Equal(saveData[:4], magic) {
		panic("invalid save file magic")
	}

	dbPath, imageStrStart, dbEndOffset, err := extractDb(saveData, args.DumpDB, args.InjectDB, args.OutPath)
	if err != nil {
		panic(err)
	}

	if args.DumpDB {
		os.Exit(0)
	}

	if !args.InjectDB {
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

		err = writeToDb(db, args)
		if err != nil {
			panic(err)
		}
	}

	updatedDbBytes, err := os.ReadFile(dbPath)
	if err != nil {
		panic(err)
	}

	if args.InjectDB {
		if !bytes.Equal(updatedDbBytes[:15], dbMagic) {
			panic("invalid db magic")
		}
		args.OutPath = args.InPath
	}

	err = writeSave(updatedDbBytes, saveData, imageStrStart, dbEndOffset, args.OutPath)
	if err != nil {
		panic(err)
	}
}