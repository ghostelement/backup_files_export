package checkfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Jobs struct {
	Jobs []Job
}
type Job struct {
	Name           string `yaml:"name"`
	FileDir        string `yaml:"fileDir"`
	FileDateFormat string `yaml:"fileDateFormat"`
	FileType       string `yaml:"fileType"`
}

type CheckList struct {
	Name  string
	Count int
	Size  int64
}

// 解析yaml文件
func LoadPath(path string) (*Jobs, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := []Job{}
	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, err
	}
	jobs := Jobs{
		Jobs: c,
	}
	return &jobs, nil

}

func (c *Job) CheckFile() (int, int64, error) {
	i := 0
	var size int64
	//生成日期格式
	if c.FileDateFormat == "" {
		c.FileDateFormat = "20060102"
	}
	backupKey := time.Now().Format(c.FileDateFormat)
	//fmt.Println("time key: ", backupKey)
	// 获取目录下文件
	files, err := os.ReadDir(c.FileDir)
	if err != nil {
		return i, size, err
	}
	// 根据文件类型和日期格式进行匹配
	for _, file := range files {
		//fmt.Println(strings.ToLower(filepath.Ext(file.Name())))
		//filetype := strings.ToLower(filepath.Ext(file.Name()))[1:]
		//fmt.Println(filepath.Ext(file.Name()))
		if file.IsDir() && c.FileType == "" {
			//fmt.Println(file.Name())
			if strings.Contains(file.Name(), backupKey) {
				i++
				continue
			}
		}
		// 获取文件类型
		filetype := filepath.Ext(file.Name())
		// 去除.
		if filetype != "" {
			filetype = strings.Replace(filetype, ".", "", -1)
			//fmt.Println(filetype)
		}
		// 文件类型为空则只匹配文件名
		if c.FileType == "" {
			filetype = ""
		}
		if strings.Contains(file.Name(), backupKey) && (filetype == c.FileType) {
			fileinfo, err := os.Stat(fmt.Sprint(c.FileDir, "/", file.Name()))
			//fmt.Println("fileinfo: ", fileinfo)
			if err != nil {
				fmt.Println(err)
			}
			if fileinfo != nil {
				size = size + fileinfo.Size()
				//fmt.Println(file.Name(), " size is : ", size)
			}
			i++
		}
	}
	return i, size, nil
}

// 过滤获取文件存在状态
func GetFileStat(path string) (*[]CheckList, error) {
	list := []CheckList{}
	jobs, err := LoadPath(path)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs.Jobs {
		c := CheckList{}
		count, size, err := job.CheckFile()
		if err != nil {
			return nil, err
		}
		//fmt.Printf("job: %s, count: %d\n", job.Name, count)
		c.Name = job.Name
		c.Count = count
		c.Size = size
		list = append(list, c)
	}
	return &list, nil
}
