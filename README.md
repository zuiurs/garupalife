# garupalife

Motto Garupa Life!

## Usage

Download Garupa 4-panel comics.

```
$ garupalife.exe scrape
```

Scraping options:

|options|description|
|:-:|:-|
|`-d` `--directory`|Output directory (Default: `.`)|
|`-t` `--time`|Scraping interval in ms (Default: `500`) |
|`--offset`|Backnumber offset (Default: `1`)|

```
$ garupalife.exe scrape -d output --offset 47 -t 1000
Extract URLs: 41 - 50
Extract URLs: 51 - 60
Extract URLs: 61 - 70
Saving Image: output/47_motto_garupa_life.jpg
Saving Image: output/48_motto_garupa_life.jpg
Saving Image: output/49_motto_garupa_life.jpg
Saving Image: output/50_motto_garupa_life.jpg
Saving Image: output/51_motto_garupa_life.jpg
Complete saving all images
```

Help.

```
$ garupalife
motto garupa life command line tool

Usage:
  garupalife [command]

Available Commands:
  help        Help about any command
  scrape      scrape four-panel comic

Flags:
  -h, --help   help for garupalife

Use "garupalife [command] --help" for more information about a command.
```

## Installation

```
go get -u github.com/zuiurs/garupalife
```

## License

This software is released under the MIT License, see LICENSE.
