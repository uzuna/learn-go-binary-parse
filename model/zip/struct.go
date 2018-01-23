package zip

/*
 CentralDirectoryの情報だけ
*/
type ZipVerbose struct {
	EndOfCentralDirectory   EndOfCentralDirectory
	CentralDirectoryHeaders []CentralDirectoryHeader
}

/*
 End of Central Directory Record
*/
type EndOfCentralDirectory struct {
	Signeture               uint32
	DiskNumber              uint16
	DiskStartNumber         uint16
	EntriesNumberOfThisDisk uint16
	EntriesNumber           uint16
	Size                    uint32
	Offset                  uint32
	CommentLength           uint16
	// Comment         uint16
}

/*
 Central Directory Header
*/
type CentralDirectoryHeader struct {
	Signeture                   uint32
	VersionMadeBy               []uint8
	VersionNeededToExtract      uint16
	GeneralPurposeBitFlag       uint16
	CompressionMethod           uint16
	LastModFileTime             uint16
	LastModFileDate             uint16
	CRC32                       uint32
	CompressedSize              uint32
	UnCompressedSize            uint32
	FileNameLength              uint16
	ExtraFieldLength            uint16
	FileCommentLength           uint16
	DiskNumberStart             uint16
	InternalFileAttributes      uint16
	ExternalFileAttributes      uint32
	RelativeOffsetOfLocalHeader uint32

	FileName    string
	ExtraField  []uint8
	FileComment []uint8

	// 終端位置を取得
	HeaderLength uint32
}

/*
 Local header
*/
type LocalFile struct {
	Signeture              uint32
	VersionNeededToExtract uint16
	GeneralPurposeBitFlag  uint16
	CompressionMethod      uint16
	LastModFileTime        uint16
	LastModFileDate        uint16
	CRC32                  uint32
	CompressedSize         uint32
	UnCompressedSize       uint32
	FileNameLength         uint16
	ExtraFieldLength       uint16

	FileName   string
	ExtraField []uint8

	FileData []byte

	HeaderLength uint32
}
