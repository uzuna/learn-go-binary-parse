## Golangでbinaryを操作する

ダミーのzipfileを使って
https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT

### Goの使用関係

package import の制約の関係からfileの配置を$GOPATH配下の`github.com/uzuna/learn-go-binary-parse`に移動。


### analysis zip file format

Big EndianはL->R
Little EndianはR->L


##### file structure

- zip local file header A
- file entrie A
- ZLFH B
- FE B
- zip central directory file header A
- zip central directory file header B
- zip central directory end code

#### CD

```
4.3.12  Central directory structure:

[central directory header 1]
.
.
. 
[central directory header n]
[digital signature] 

File header:

central file header signature   4 bytes  (0x02014b50)
version made by                 2 bytes
version needed to extract       2 bytes
general purpose bit flag        2 bytes
compression method              2 bytes
last mod file time              2 bytes
last mod file date              2 bytes
crc-32                          4 bytes
compressed size                 4 bytes
uncompressed size               4 bytes
file name length                2 bytes
extra field length              2 bytes
file comment length             2 bytes
disk number start               2 bytes
internal file attributes        2 bytes
external file attributes        4 bytes
relative offset of local header 4 bytes

file name (variable size)
extra field (variable size)
file comment (variable size)

 4.3.13 Digital signature:

        header signature                4 bytes  (0x05054b50)
        size of data                    2 bytes
        signature data (variable size)
```

実データ 

```
50 4B 01 02 // S
14 00 // MadeBy
14 00 // NeedExtract
08 00 // general purpose
08 00 // Compression
00 00 // time
00 00 // date
2D 73 07 F0 // CRC
09 00 00 00 // C-size
03 00 00 00 // UC-size
0B 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 64 75 6D 6D 79 2F 61 2E 74 78 74 
```


```
50 4B 01 02 
14 00 
14 00 
00 00 /GP
08 00 /Comp
FC 66 //time
35 4C //date
48 02 07 39 05 00 00 00 06 00 00 00 05 00 00
00 00 00 00 00 01 00 20 00 00 00 26 00 00 00 62
2E 74 78 74
```

#### zip central directory end code

zip end of central dir 
```
4.3.16  End of central directory record:

end of central dir signature    4 bytes  (0x06054b50)
number of this disk             2 bytes
number of the disk with the
start of the central directory  2 bytes
total number of entries in the
central directory on this disk  2 bytes
total number of entries in
the central directory           2 bytes
size of the central directory   4 bytes
offset of start of central
directory with respect to
the starting disk number        4 bytes
.ZIP file comment length        2 bytes
.ZIP file comment       (variable size)
```

zip central dhirectory end codeの始点`50 4B 05 06`
