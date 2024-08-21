# Yew Programming Language

Welcome to the Yew language repo!

Official Yew language site (not yet up as of 08/21/2024):
- yew-lang.org

## Installing
TODO

### Building from Source
TODO

## Commands
TODO

### Command `yew`: Yew Compiler

|---|---|---|---|
|Command|Option(s)|Description|Example|
| | | | |
| | |Starts repl interface (see `repl` command)|`yew`|
| | | | |
|`repl`|| | |
| |*none*|Starts repl interface|`yew repl`|
| |`-i [pkg1,pkg2,..]`|Imports pkg1, pkg2, ...|`yew repl -i base,reflect`|
| |&emsp;&emsp;`--import`| | |
| |`-L`|Runs in literate mode|`yew repl -L`| |
| |&emsp;&emsp;`--literate`| | |
| |`-o [file]`|Outputs REPL input to `file`|`yew repl -o record.yew`|
| |&emsp;&emsp;`--out`| | |
| |&emsp;&emsp;`--output`| | |
| | | | |
|`build`| | | |
| |*none*|Builds package found in pwd|`yew build`|
| |`[pkg]`|Builds `pkg`. Must be first arg.|`yew build pkg`|
| |`-o <name>`|Names executable output `name`|`yew build pkg -o a.out`|
| |&emsp;&emsp;`--out`| | |
| |&emsp;&emsp;`--output`| | |
| |`-- <pkg>`|Builds package `pkg`|`yew -o a.out -- pkg`|
| |`-i`|Stops compilation after producing all IR|`yew pkg -i`|
| |&emsp;&emsp;`--ir`|||
| |&emsp;&emsp;`--intermediate`|||
| |&emsp;&emsp;`--Intermediate`|||
| |`-w all`|Enables all warnings|`yew build -w all`|
| |&emsp;&emsp;`--warn all`|||
| |&emsp;&emsp;`--warning all`|||
| |`-w none`|Disables all warnings|`yew build -w none`|
| |&emsp;&emsp;`--warn none`|||
| |&emsp;&emsp;`--warning none`|||
| |`-w <config>`|Uses warning flags described in `config`|`yew build -w warn.config`|
| |&emsp;&emsp;`--warn`|||
| |&emsp;&emsp;`--warning`|||
| | | | |
|`help`| | | |
| |*none*|Displays info for common commands|`yew help`|
| |`[cmd]`|Displays info for `cmd`. Must be first arg.|`yew help build`|
| |`-o <option>`|Displays info for `option` of cmd|`yew help build -o ir`|
| |&emsp;&emsp;`--opt`|||
| |&emsp;&emsp;`--option`|||
| |`-- <cmd>`|Displays info for command `cmd`|`yew help -o ir -- build`|
| |`-v [bool]`|Sets verbose help to `bool` (def. `true`)|`yew help build -v`|
| |&emsp;&emsp;`--verbose`|||
| | | | |
|`version`| | | |
| |*none*|Displays running version of yew compiler|`yew version`|
| | | | |
|`...`|`...`|...|`...`|

TODO: finish

### Command: `ypk`: Yew Package Manager 
TODO

## License
Yew is distributed under the terms of the MIT license.

See the file `LICENSE` located in the same directory as this file for more details