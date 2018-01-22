
/*
 Directoryの構造体
*/
type Directory struct {
	path    string
	entries []interface{} // 配下にDirectoryorFileを持つ
}

/*
 Fileの構造体
*/
type File struct {
	name string
	body []byte
}