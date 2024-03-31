package merkledag

import (
	"encoding/json"
	"strings"
)

const STEP = 4

// Hash to file

func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据hash和path， 返curObjBinary回对应的文件, hash对应的类型是tree
	// 查询Hash是否存在
	flag, _ := store.Has(hash)
	if flag {
		//根据给定的键从存储中检索对应的值
		objBinary, _ := store.Get(hash)
		obj := binaryToObj(objBinary)
		pathArr := strings.Split(path, "\\")
		cur := 1
		return getFileByDir(obj, pathArr, cur, store)
	}
	return nil
}

// 根据path， 返回对应的文件
func getFileByDir(obj *Object, pathArr []string, cur int, store KVStore) []byte {
	if cur >= len(pathArr) {
		return nil
	}
	index := 0
	for i := range obj.Links {
		// 获取存储对象的类型
		objType := string(obj.Data[index : index+STEP])
		index += STEP
		objInfo := obj.Links[i]
		//查找跟传入的路径名字相同的obj
		if objInfo.Name != pathArr[cur] {
			continue
		}
		//匹配对应的类型
		switch objType {
		case TREE:
			objDirBinary, _ := store.Get(objInfo.Hash)
			objDir := binaryToObj(objDirBinary)
			//如果类型为tree，递归
			ans := getFileByDir(objDir, pathArr, cur+1, store)
			if ans != nil {
				return ans
			}
		case BLOB:
			ans, _ := store.Get(objInfo.Hash)
			return ans
		case LIST:
			objLinkBinary, _ := store.Get(objInfo.Hash)
			objList := binaryToObj(objLinkBinary)
			ans := getFileByList(objList, store)
			return ans
		}
	}
	return nil
}

func getFileByList(obj *Object, store KVStore) []byte {
	ans := make([]byte, 0)
	index := 0
	for i := range obj.Links {
		curObjType := string(obj.Data[index : index+STEP])
		index += STEP
		curObjLink := obj.Links[i]
		curObjBinary, _ := store.Get(curObjLink.Hash)
		curObj := binaryToObj(curObjBinary)
		if curObjType == BLOB {
			ans = append(ans, curObjBinary...)
		} else { //Type为List
			tmp := getFileByList(curObj, store)
			ans = append(ans, tmp...)
		}
	}
	return ans
}

func binaryToObj(objBinary []byte) *Object {
	var res Object
	//将json转换为结构体Object
	json.Unmarshal(objBinary, &res)
	return &res
}
