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
Use this with the unstuck flag to prevent this or just fast travel.

Customise wand:   
`hlse_x64.exe -i HL-00-00.sav --wand-base soft_spiral_warm_brown --wand-core dragon_heartstring`

Dump DB:   
`hlse_x64.exe -i HL-00-00.sav -o out.db --dump-db`

Inject DB:   
`hlse_x64.exe -i out.db -o HL-00-00.sav --inject-db`

Probe save:    
`hlse_x64.exe -i HL-00-00.sav -p`
```
Name:           Eve -
XP:             2800
Level:          5
House:          Hufflepuff
Galleons:       737
Inventory size: 20
Talent points:  1

Wand:
Base:   Soft spiral, light brown
Wood:   Fir
Core:   Dragon heartstring
Length: 12 ½
Flex:   Reasonably supple

Resource inventory:
Knuts, 737
LacewingFlies, 3
LeapingToadstool_Byproduct, 1
Moonstone, 5
```

```
Usage: hlse_x64.exe --inpath INPATH [--outpath OUTPATH] [--xp XP] [--level LEVEL] [--galleons GALLEONS] [--first-name FIRST-NAME] [--last-name LAST-NAME] [--house HOUSE] [--inventory-size INVENTORY-SIZE] [--probe] [--item-quantities ITEM-QUANTITIES] [--talent-points TALENT-POINTS] [--unstuck] [--reset-talent-points] [--dump-db] [--inject-db] [--wand-base WAND-BASE] [--wand-wood WAND-WOOD] [--wand-core WAND-CORE] [--wand-length WAND-LENGTH] [--wand-flex WAND-FLEX]

Options:
  --inpath INPATH, -i INPATH
                         Path of input save file.
  --outpath OUTPATH, -o OUTPATH
                         Path of output save file.
  --xp XP                Set XP.
  --level LEVEL          Set player's level.
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
  --wand-base WAND-BASE
                         Set wand base.
  --wand-wood WAND-WOOD
                         Set wand wood.
  --wand-core WAND-CORE
                         Set wand core.
  --wand-length WAND-LENGTH
                         Set wand length.
  --wand-flex WAND-FLEX
                         Set wand flex.
  --help, -h             display this help and exit
```

## Wand options
### bases
```
notched_warm_brown
notched_light_brown
notched_dusty_pink
classic_grey
classic_gray
classic_black
classic_grey-brown
classic_gray-brown
classic_grey_brown
classic_gray_brown
soft_spiral_light_brown
soft_spiral_warm_brown
soft_spiral_light_black
spiral_ash_brown
spiral_green-grey
spiral_green-gray
spiral_green_grey
spiral_green_gray		
spiral_dark_brown
stalk_honey_brown
stalk_dark_brown		
stalk_warm_brown
ringed_dark_brown
ringed_pale_brown
ringed_buff
crooked_spiral_dark_grey
crooked_spiral_dark_gray
crooked_spiral_warm_brown
crooked_spiral_pale_brown
natural_grey
natural_gray
natural_honey_brown
natural_warm_brown
```

### woods
```
acacia
alder
ash
beech
blackthorn
black_walnut
cedar
cherry
chestnut
cypress
dogwood
ebony
elder
elm
english_oak
fir
hawthorn
holly
hazel
hornbeam
larch
laurel
maple				
pear
pine
poplar
red_oak
redwood
rowan
silver_lime
spruce
sycamore
vine
walnut
willow
yew
```

### cores
```
dragon_heartstring
phoenix_feather
unicorn_hair
```

### lengths
```
nine_and_a_half
nine_and_three_quarters
ten
ten_and_a_quarter
ten_and_a_half
ten_and_three_quarters
eleven
eleven_and_a_quarter
eleven_and_a_half
eleven_and_three_quarters
twelve
twelve_and_a_quarter
twelve_and_a_half
twelve_and_three_quarters
thirteen	
thirteen_and_a_quarter
thirteen_and_a_half
thirteen_and_three_quarters
fourteen
fourteen_and_a_quarter
fourteen_and_a_half
```

### flex
```
brittle
fairly_bendy
hard
pliant
quite_bendy
quite_flexible
reasonably_supple
rigid
slightly_springy
slightly_yielding
solid
stiff
supple
surprisingly_swishy
swishy
unbending
unyielding
very_flexible
```

## Disclaimer
- I will not be responsible for any possibility of save corruption.    
- Hogwarts Legacy brand and name is the registered trademark of its respective owner.    
- HLSE has no partnership, sponsorship or endorsement with Avalanche Software or Warner Bros. Games.
