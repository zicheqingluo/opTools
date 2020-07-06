package main
import (
	"encoding/json"
	"os/exec"
	"fmt"
	"strings"
	"bytes"
	"log"
	"io/ioutil"
	"flag"
)
var fileMap map[string]string
var checkFlag bool

func getInode(execDir string) map[string]string{
	cmd := exec.Command("ls","-i",execDir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	outStr:=out.String()
	fileMap=make(map[string]string)
	for _,line := range strings.Split(outStr,"\n"){
		if len(line) >0 {

			s := strings.Split(line," ")
			fileMap[s[0]]=s[1]
		}
		
	}
	//fmt.Println(fileMap)
	return fileMap
}

func diffChanges(f map[string]string){
	content, err := ioutil.ReadFile("./fdata")
	if err != nil {
		fmt.Println("read file failed, err:", err)
		return
	}
	var fj	map[string]string 
	err = json.Unmarshal(content,&fj)
	if err != nil {
		fmt.Println("unmarshal inode data failed,err:",err)
	}
	for k := range f{
		_,ok := fj[k]
		if !ok {
			fmt.Println("新增文件:",k,f[k])
			checkFlag = true
		}
	}

	// for k := range fj{
	// 	_,ok := f[k]
	// 	if !ok {
	// 		delete(fj, k)
	// 	}
	// }
}

func fileSave(f map[string]string){
	fj,err := json.Marshal(f)
	if err != nil {
		fmt.Printf("inode marshal fail:%v",err)
		return 
	}
	err = ioutil.WriteFile("./fdata", fj, 0644)
	if err != nil {
		fmt.Println("write file failed, err:", err)
		return
	}

}

func main(){
	//1、通过flag动态获取要检测的目录
	//2、生成该目录的inode 和文件名的对应关系，序列化为json
	//3、对比
		//4.1 读取上一次持久化到本地的记录
		//4.2 遍历本次记录与上一次做对比，看是否有新增
		//4.3 遍历上一次记录与本次做对比，看是否有文件删除
	//4、将对应关系持久化到本地
	//5、返回是否有变化
	var execDir string
	flag.StringVar(&execDir, "dir", "/home/yxk/test", "目标目录")
	flag.Parse()

	fileMap:=getInode(execDir)
	diffChanges(fileMap)
	fileSave(fileMap)
	if checkFlag{
		fmt.Printf("coredump:0")
	}
}