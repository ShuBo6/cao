package main

import (
	"fmt"
	"gopkg.in/resty.v1"
	"os"
	"path"
	"strings"
	"sync"
)

func main() {
	//dfsPath("data/我的阿里云/极客时间/01-专栏课")
	downloadOneFile("data/我的阿里云/极客时间/01-专栏课/001-050/01-数据结构与算法之美/01-开篇词 (1讲)/url")
}

func dfsPath(rootPath string) {
	dir, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range dir {
		currPath := path.Join(rootPath, entry.Name())
		if entry.IsDir() {
			dfsPath(currPath)
		} else {
			fmt.Println(currPath)
			downloadOneFile(currPath)
		}

	}
}

func downloadOneFile(filepath string) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(b), "\n")
	wg := sync.WaitGroup{}
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lines[i] = strings.TrimRight(lines[i], "\r")
			r := resty.New()
			fmt.Println("下载url为", lines[i])
			resp, err := r.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).R().Get(lines[i])
			if err != nil {
				fmt.Println(err)
				return
			}
			err = os.WriteFile(path.Join(path.Dir(filepath), path.Base(lines[i])), resp.Body(), 0777)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(i)
	}
	wg.Wait()
}
