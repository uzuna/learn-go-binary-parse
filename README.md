## Golangでbinaryを操作する

ダミーのzipfileを使って



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

#### zip central directory end code

zip64 end of central dir 
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
