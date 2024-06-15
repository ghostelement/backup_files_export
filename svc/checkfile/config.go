package checkfile

import (
	"os"
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

func (c *Job) CheckFile() (int, error) {
	i := 0
	//生成日期格式
	if c.FileDateFormat == "" {
		c.FileDateFormat = "20060102"
	}
	backupKey := time.Now().Format(c.FileDateFormat)
	//fmt.Println("time key: ", backupKey)
	// 获取目录下文件
	files, err := os.ReadDir(c.FileDir)
	if err != nil {
		return i, err
	}
	// 根据文件类型和日期格式进行匹配
	for _, file := range files {
		//fmt.Println(strings.ToLower(filepath.Ext(file.Name())))
		//filetype := strings.ToLower(filepath.Ext(file.Name()))[1:]
		if strings.Contains(file.Name(), backupKey) && strings.Contains(file.Name(), c.FileType) {
			i++
		}
	}
	return i, nil
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
		count, err := job.CheckFile()
		if err != nil {
			return nil, err
		}
		//fmt.Printf("job: %s, count: %d\n", job.Name, count)
		c.Name = job.Name
		c.Count = count
		list = append(list, c)
	}
	return &list, nil
}
