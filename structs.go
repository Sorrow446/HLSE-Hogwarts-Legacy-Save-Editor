package main

type Args struct {
	InPath           string   `arg:"-i, required" help:"Path of input save file."`
	OutPath          string   `arg:"-o" help:"Path of output save file."`
	XP               int      `arg:"--xp" help:"Set XP."`
	Galleons         int      `arg:"--galleons" help:"Set Galleons."`
	FirstName        string   `arg:"--first-name" help:"Set character first name."`
	Surname          string   `arg:"--last-name" help:"Set character last name."`
	InventorySize    int      `arg:"--inventory-size" help:"Set inventory size."`
	Probe            bool     `arg:"-p" help:"Probe save file and exit."`
	ItemQuantities   []string `arg:"--item-quantities" help:"Set quantities of inventory items. <item ID> <quantity> pairs separated by spaces."`
	ParsedItemQuants []*ItemPairs `arg:"-"`
	TalentPoints     int      `arg:"--talent-points" help:"Set talent points."`
}

type ItemPairs struct {
	ItemID   string
	Quantity int64
}