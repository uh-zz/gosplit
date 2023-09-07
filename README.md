# gosplit

## Name

gosplit - split a file into pieces in Go

## Description

gosplit is a like Unix command that uses only the Go standard library.


## Options

### `-b (byte count)`

`-b` option splits the file by specified bytes

#### Example

```
❯ echo aaa | ./gosplit -b 1

❯ ls
xaa  xab  xac  xad
```

### `-l (line count)`

`-l` option splits the file by specified lines


#### Example
```
❯ echo aaa | ./gosplit -l 1

❯ ls
xaa
```
