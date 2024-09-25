package utils

import (
	"fmt"
	"io"
	"os"
	"path"
)

// CopyFile 复制文件。
// srcName 是源文件路径。
// dstName 是目标文件路径。
// 返回值 written 是写入目标文件的字节数。
// 返回值 err 是复制过程中可能发生的错误。
func CopyFile(srcName string, dstName string) (written int64, err error) {
	// 打开源文件以读取。
	src, err := os.Open(srcName)
	if err != nil {
		// 如果无法打开源文件，返回错误。
		return
	}
	// 打开或创建目标文件以写入。
	// 使用0644权限确保文件所有者可读写，其他用户可读。
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		// 如果无法打开或创建目标文件，返回错误。
		return
	}

	// 确保在函数结束前关闭文件。
	defer func() {
		if src != nil {
			err := src.Close()
			if err != nil {
				return
			}
		}
		if dst != nil {
			_ = dst.Close()
		}
	}()
	// 将源文件内容复制到目标文件。
	// 返回复制的字节数和可能发生的错误。
	return io.Copy(dst, src)
}

// FileTmpPath 生成一个临时文件的路径。
// 该函数接收一个完整的文件路径（fullPath）和一个数据库前缀（dbPrefix）作为参数，
// 并返回一个基于这些参数生成的临时文件路径。
// 参数：
// - fullPath: 完整的文件路径，用于计算临时文件的唯一标识。
// - dbPrefix: 数据库前缀，用于区分不同数据库源的临时文件。
// 返回值：
// - string: 生成的临时文件路径。
func FileTmpPath(fullPath, dbPrefix string) string {
	// 获取文件的后缀名。
	suffix := path.Ext(fullPath)
	// 检查后缀名长度，如果为0，则记录错误。
	// 这是因为有效的文件路径应该包含文件后缀，否则可能无法正确处理文件。
	if len(suffix) == 0 {
		fmt.Println("suffix  err:")
	}

	// 返回临时文件路径，由数据库前缀、原文件路径的MD5值和原文件后缀组成。
	// 这样做可以确保临时文件名的唯一性，同时保留原文件的类型。
	return dbPrefix + Md5(fullPath) + suffix //a->b
}

// FileExist 检查文件是否存在
// 该函数通过os.Stat检查文件是否存在，而不是直接打开文件，
// 这样可以避免不必要的文件打开操作，提高效率。
// 参数:
//
//	filename: 需要检查的文件名
//
// 返回值:
//
//	bool: 文件是否存在
func FileExist(filename string) bool {
	// os.Stat会返回文件的信息，如果文件不存在，则返回错误
	// 我们通过检查错误类型来判断文件是否存在
	_, err := os.Stat(filename)
	// 当err为nil时，表示文件存在；当err不为nil但错误类型为os.IsExist时，也表示文件存在
	return err == nil || os.IsExist(err)
}
