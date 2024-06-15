# backup_files_export
## 项目介绍
本项目是一个简单的备份文件检查exporter监控，适用于prometheus监控。  
其主要功能是帮助检查指定目录下的备份文件是否存在，并生成监控指标暴露给prometheus。  

## 用法
通过指定的配置文件启动程序，暴露prometheus的监控指标。  
```sh
backup_files_export --metrics.config=config.yaml --metrics.port=9103
```
config.yaml配置文件示例：  
```yaml
- name: mysql
  # 备份文件目录
  fileDir: "/backup/mysql"
  # 备份文件日期格式
  fileDateFormat: "20060102"
  # 备份文件类型
  fileType: "sql" 
- name: tidb
  fileDir: "/backup/tidb"
  fileDateFormat: "20060102"
  fileType: ""
```
