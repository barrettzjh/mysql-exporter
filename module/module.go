package module

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"

	//"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	c = Conf{}
	Configure = c.getConf()
	Client = GetMysqlClient()
)

type Conf struct {
	Port     string `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
	MysqlUrl string `yaml:"mysqlUrl"`
}

func (c *Conf) getConf() *Conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func GetMysqlClient()map[string]*sql.DB{
	var db *sql.DB
	var err error
	sliceClient := make(map[string]*sql.DB)
	var index string
	for _, v := range strings.Split(Configure.MysqlUrl, ";"){
		index = strings.Split(strings.Split(v, "tcp(")[1], ":")[0]
		db, err = sql.Open("mysql", v)
		if err != nil {
			fmt.Println(err)
		}
		sliceClient[index] = db
	}
	return sliceClient
}

func GetMysqlStorage()map[string]map[string]string{
	var schema, gb string
	mapStorage := make(map[string]string)
	sliceStorage := make(map[string]map[string]string)
	for f,j := range Client{
		mapStorage = map[string]string{}
		rows, err := j.Query("SELECT table_schema,SUM(AVG_ROW_LENGTH*TABLE_ROWS+INDEX_LENGTH) AS total_mb FROM information_schema.TABLES group by table_schema;\n")
		if err != nil{
			fmt.Println(err)
		}
		for rows.Next(){
			err := rows.Scan(&schema, &gb)
			if schema == "information_schema"{
				continue
			}
			if err != nil{
				fmt.Println(err)
			}
			mapStorage[schema] = gb
		}
		sliceStorage[f] = mapStorage
	}
	return sliceStorage
}

func GetMysqlStatus()map[string]map[string]string{
	var name, value string
	mapStatus := make(map[string]string)
	sliceStatus := make(map[string]map[string]string)
	for f,j := range Client{
		mapStatus = map[string]string{}
		rows, err := j.Query("show status;")
		if err != nil{
			fmt.Println(err)
		}
		for rows.Next(){
			err = rows.Scan(&name, &value)
			if err != nil{
				fmt.Println(err)
			}
			mapStatus[name] = value
		}
		sliceStatus[f] = mapStatus
	}
	return sliceStatus
}

func StringToFloat(v string) (float64, error){
	v2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}
	return v2, err
}