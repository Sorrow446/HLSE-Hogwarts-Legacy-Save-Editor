# HLSE - Hogwarts-Legacy-Save-Editor
Save editor for Hogwarts Legacy written in Go.    
[Windows binaries](https://github.com/Sorrow446/HLSE-Hogwarts-Legacy-Save-Editor/releases)

## Usage
Set inventory size to 100 and galleons to 1000 and overwrite:   
`hlse_x64.exe -i HL-00-00.sav --inventory-size 100 --galleons 1000`

Set first and last name and write to out.sav:   
`hlse_x64.exe -i HL-00-00.sav --first-name "first name" --last-name "last name" -o out.sav`

Set quantities of inventory items:   
`hlse_x64.exe -i HL-00-00.sav --item-quantities StenchOfTheDead 100 LacewingFlies 10`

Set player's house to Hufflepuff:   
`hlse_x64.exe -i HL-00-00.sav --house hufflepuff`

Be careful. If the player's currently in another house's common room at the time, the player will get locked in there.
Use this with the unstuck flag to prevent this or fast travel.

Dump DB:   
`hlse_x64.exe -i HL-00-00.sav -o out.db --dump-db`

Inject DB:   
`hlse_x64.exe -i out.db -o HL-00-00.sav --inject-db`

Probe save:    
`hlse_x64.exe -i HL-00-00.sav -p`
```
Name:           first name last name
XP:             2000
House:          Hufflepuff
Galleons:       1000
Inventory size: 100

Inventory:
Knuts, 1000
LacewingFlies, 84
LeapingToadstool_Byproduct, 120
Moonstone, 654
AshwinderEggs, 37
```

```
Usage: hlse_x64.exe --inpath INPATH [--outpath OUTPATH] [--xp XP] [--galleons GALLEONS] [--first-name FIRST-NAME] [--last-name LAST-NAME] [--house HOUSE] [--inventory-size INVENTORY-SIZE] [--probe] [--item-quantities ITEM-QUANTITIES] [--talent-points TALENT-POINTS] [--unstuck] [--reset-talent-points] [--dump-db] [--inject-db]

Options:
  --inpath INPATH, -i INPATH
                         Path of input save file.
  --outpath OUTPATH, -o OUTPATH
                         Path of output save file.
  --xp XP                Set XP.
  --galleons GALLEONS    Set Galleons.
  --first-name FIRST-NAME
                         Set character first name.
  --last-name LAST-NAME
                         Set character last name.
  --house HOUSE          Set player's house. gryffindor, ravenclaw, hufflepuff, slytherin
  --inventory-size INVENTORY-SIZE
                         Set inventory size.
  --probe, -p            Probe save file and exit.
  --item-quantities ITEM-QUANTITIES
                         Set quantities of inventory items. <item ID> <quantity> pairs separated by spaces.
  --talent-points TALENT-POINTS
                         Set talent points.
  --unstuck              Sets the player's coordinates to The Great Hall.
  --reset-talent-points
                         Resets the player's talent points.
  --dump-db              Dump DB from save file and exit.
  --inject-db            Inject DB into save file.
  --help, -h             display this help and exit
```

## Disclaimer
- I will not be responsible for any possibility of save corruption.    
- Hogwarts Legacy brand and name is the registered trademark of its respective owner.    
- HLSE has no partnership, sponsorship or endorsement with Avalanche Software or Warner Bros. Games.
