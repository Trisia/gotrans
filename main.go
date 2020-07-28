package main

import "io/ioutil"

func main() {

	// 繁体转换简体
	//cc, err := opencc.NewOpenCC("t2s")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//src, err := os.Open("C:\\Users\\pc\\Desktop\\135\\OEBPS\\Text\\Section0006.xhtml")
	//dst, err := os.Create("C:\\Users\\pc\\Desktop\\135\\OEBPS\\Text\\Section0006s.xhtml")
	//if err != nil {
	//	panic(err)
	//}
	//err = cc.ConvertFile(src, dst)
	//nText,err := cc.ConvertText(`繁體到简体`)
	//if err != nil{
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(nText)

	ioutil.TempDir("", "translaor")
}
