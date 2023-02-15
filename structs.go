package main

type Args struct {
	InPath            string       `arg:"-i, required" help:"Path of input save file."`
	OutPath           string       `arg:"-o" help:"Path of output save file."`
	XP                int          `arg:"--xp" help:"Set XP."`
	Level             int          `arg:"--level" help:"Set player's level."`
	Galleons          int64        `arg:"--galleons" help:"Set Galleons."`
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
	WandBase          string       `arg:"--wand-base" help:"Set wand base."`
	WandWood          string       `arg:"--wand-wood" help:"Set wand wood."`
	WandCore          string       `arg:"--wand-core" help:"Set wand core."`
	WandLength        string       `arg:"--wand-length" help:"Set wand length."`
	WandFlex          string       `arg:"--wand-flex" help:"Set wand flex."`
	Wand  			  *Wand        `arg:"-"`
}

type ItemPairs struct {
	ItemID   string
	Quantity int64
}

type Wand struct {
	Base   string
	Wood   string
	Core   string
	Length string
	Flex   string
	Style  string
}