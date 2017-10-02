# escrabble
Scrabble Piece Generator

## Piece Generation

Create a hamming code for letters in Spanish given the frequencies
at which they appear in the Spanish dictionary.

Assign pieces values based on the length of the code

The rarer the letter, the longer the code, the more points that piece is worth.

This all assumes that the rarer the letter the harder it is to use in a word.

## Dictionary Corpus

Spanish.dic was obtained from https://github.com/titoBouzout/Dictionaries/blob/master/Spanish.dic

## Running
`./run.sh`
OR
```
go build main.go
./main.go {dic_file_name}
python3 pieces.py
```

## Purpose
This prints out a png meant to be printed, laser-cut and created into
actual scrabble pieces
