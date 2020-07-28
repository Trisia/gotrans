package main

import (
	"fmt"
	"github.com/sgoby/opencc"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//for idx, arg := range os.Args {
	//	fmt.Printf("Arg [%d]: %s\n", idx, arg)
	//}
	//time.Sleep(time.Second * 5)
	//if len(os.Args) < 1 {
	//	log.Fatalf("%s [epub_file_path], 请拖拽文件到exe程序之上", os.Args[0])
	//	return
	//}
	//
	//epubFilePath := os.Args[1]

	epubFilePath := "C:\\Users\\pc\\Desktop\\OVERLORD地14卷滅國的魔女.zip"

	// 创建用于解压的工作目录
	workDir, err := ioutil.TempDir("", "gotrans")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// 删除工作目录
		os.RemoveAll(workDir)
	}()
	// 解压epub文件
	err = Unzip(epubFilePath, workDir)
	if err != nil {
		log.Fatal(err)
	}

	// 繁体转换简体
	cc, err := opencc.NewOpenCC("t2s")
	if err != nil {
		log.Fatal(err)
	}
	// 需要翻译的目录
	basePath := filepath.Join(workDir, "OEBPS", "Text")

	// 读取目录信息
	dir, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}
	// 遍历需要翻译的目录
	for _, file := range dir {
		// 过滤 xhtml 格式
		ok := strings.HasSuffix(file.Name(), ".xhtml")
		if !ok {
			continue
		}
		// 拼接出完整路径
		srcFilePath := filepath.Join(basePath, file.Name())
		srcFile, err := os.Open(srcFilePath)
		if err != nil {
			log.Fatal(err)
		}
		outFile, err := ioutil.TempFile(basePath, file.Name())
		if err != nil {
			log.Fatal(err)
		}
		// 繁体转换简体
		err = cc.ConvertFile(srcFile, outFile)
		if err != nil {
			log.Fatal(err)
		}
		srcFile.Close()
		outFile.Close()
		// 删除原文文件
		err = os.Remove(srcFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		// 重命名翻译后文件
		err = os.Rename(outFile.Name(), srcFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 构造新的文件名
	dstFilePath := newFileName(epubFilePath)
	// 重新压缩
	err = Zip(workDir, dstFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

// 创建一个新以“-简体”结尾的文件名
// epubFilePath: 文件路径
func newFileName(epubFilePath string) string {
	parentPath := filepath.Dir(epubFilePath)
	filename := filepath.Base(epubFilePath)
	extension := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(extension)]
	newName := fmt.Sprintf("%s-简体%s", name, extension)
	return filepath.Join(parentPath, newName)
}
