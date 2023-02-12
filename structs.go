package main

type Args struct {
	InPath            string       `arg:"-i, required" help:"Path of input save file."`
	OutPath           string       `arg:"-o" help:"Path of output save file."`
	XP                int          `arg:"--xp" help:"Set XP."`
	Galleons          int          `arg:"--galleons" help:"Set Galleons."`
	FirstName         string       `arg:"--first-name" help:"Set character first name."`
	Surname           string       `arg:"--last-name" help:"Set character last name."`
	House             string       `arg:"--house" help:"Set player's house. gryffindor, ravenclaw, hufflepuff, slytherin"`
	InventorySize     int          `arg:"--inventory-size" help:"Set inventory size."`
	Probe             bool         `arg:"-p" help:"Probe save file and exit."`
	ItemQuantities    []string     `arg:"--item-quantities" help:"Set quantities of inventory items. <item ID> <quantity> pairs separated by spaces."`
	ParsedItemQuants  []*ItemPairs `arg:"-"`
	TalentPoints      int          `arg:"--talent-points" help:"Set talent points."`
	Unstuck           bool         `arg:"--unstuck" help:"Sets the player's coordinates to The Great Hall."`
	ResetTalentPoints bool         `arg:"--reset-talent-points" help:"Resets the player's talent points."`
	DumpDB            bool         `arg:"--dump-db" help:"Dump DB from save file and exit."`
	InjectDB          bool         `arg:"--inject-db" help:"Inject DB into save file."`
}

type ItemPairs struct {
	ItemID   string
	Quantity int64
}