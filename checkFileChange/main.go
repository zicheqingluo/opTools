package main
import (
	"encoding/json"
	"os/exec"
	"os"
	"path"
	"fmt"
	"strings"
	"bytes"
	"log"
	"io/ioutil"
	"flag"
	"path/filepath"
)

var execDir,dataCache string
var fileMap map[string]string
var checkFlag bool

func getDir(){
	ex, err := os.Executable()
    if err != nil {
        panic(err)
	}
	cwdPath := filepath.Dir(ex)
	
	flag.StringVar(&execDir, "dir", "/home/homework/coresave", "目标目录")
	flag.Parse()
	name:=path.Clean(execDir)
	_,name=path.Split(name)
	dataCache=fmt.Sprintf("%s/.%s",cwdPath,name)
}

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
			l:= strings.TrimSpace(line)
			s:= strings.Split(l," ")
			//fmt.Println(s)
			fileMap[s[0]]=s[1]
		}
	}
	//fmt.Println(fileMap)
	return fileMap
}

func diffChanges(f map[string]string,dataCache string){
	content, err := ioutil.ReadFile(dataCache)
	if err != nil {
		//fmt.Println("read file failed, err:", err)
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

func fileSave(f map[string]string,dataCache string){
	fj,err := json.Marshal(f)
	if err != nil {
		fmt.Printf("inode marshal fail:%v",err)
		return 
	}
	err = ioutil.WriteFile(dataCache, fj, 0644)
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
	//_, filename, _, ok := runtime.Caller(0)

	getDir()

	fileMap:=getInode(execDir)
	diffChanges(fileMap,dataCache)
	fileSave(fileMap,dataCache)
	if checkFlag{
		fmt.Println("coredump:0")
	} else {
		fmt.Println("coredump:1")
	}
}
