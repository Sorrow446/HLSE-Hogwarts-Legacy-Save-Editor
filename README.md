# HLSE - Hogwarts-Legacy-Save-Editor
Save editor for Hogwarts Legacy written in Go.    
[Windows binaries](https://github.com/Sorrow446/HLSE-Hogwarts-Legacy-Save-Editor/releases)

## Usage
Set inventory size to 100 and galleons to 1000 and overwrite:   
`hlse_x64.exe -i HL-00-00.sav --inventory-size 100 --galleons 1000`

Set first and last name and write to out.sav:   
`hlse_x64.exe -i HL-00-00.sav --first-name "first name" --last-name "last name" -o out.sav`

Probe save:
`hlse_x64.exe -i HL-00-00.sav -p`
```
Name:           first name last name
XP:             2000
House:          Hufflepuff
Galleons:       1000
Inventory size: 100
```

```
Usage: hlse_x64.exe --inpath INPATH [--outpath OUTPATH] [--xp XP] [--galleons GALLEONS] [--first-name FIRST-NAME] [--last-name LAST-NAME] [--inventory-size INVENTORY-SIZE] [--probe]

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
  --inventory-size INVENTORY-SIZE
                         Set inventory size.
  --probe, -p            Probe save file and exit.
  --help, -h             display this help and exit
  ```

## Disclaimer
- I will not be responsible for any possibility of save corruption.    
- Hogwarts Legacy brand and name is the registered trademark of its respective owner.    
- HLSE has no partnership, sponsorship or endorsement with Avalanche Software or Warner Bros. Games.
