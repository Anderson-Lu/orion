package utils

var (
	_IPUtil     = &IPUtil{}
	_StringUtil = &StringsUtil{}
	_ArrayUtil  = &ArrayUtil{}
	_FileUtil   = &FileUtil{}
)

func IP() *IPUtil {
	return _IPUtil
}

func String() *StringsUtil {
	return _StringUtil
}

func Array() *ArrayUtil {
	return _ArrayUtil
}

func File() *FileUtil {
	return _FileUtil
}
