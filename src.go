package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	"gopkg.in/yaml.v3"
)

var global_path = os.Args[1] //conf.Arguments.Name
var _flag bool

type Hash_line struct {
	name string
	md5  string
	sha1 string
}

type Config struct {
	Arguments struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	} `yaml:"arguments"`
}

// MD5 hashes using md5 algorithm
func hash_file_md5(filePath string) string {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String
	}
	//Tell the program to call the following function when the current function returns
	defer file.Close()
	//Open a new SHA1 hash interface to write to
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String
}

// SHA1 hashes using sha1 algorithm
func hash_file_sha1(filePath string) string {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSHA1String string
	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String
	}
	//Tell the program to call the following function when the current function returns
	defer file.Close()
	//Open a new SHA1 hash interface to write to
	hash := sha1.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String
	}
	//Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]
	//Convert the bytes to a string
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String
}

// Insert all files recursivly
func FilePathWalkDir(root string) []string {
	var files []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// Another version for open blacklist to a list of hashes
func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func GetLengthFile(file string) int {
	openfile, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return 0
	}
	filedata, err := csv.NewReader(openfile).ReadAll()
	if err != nil {
		log.Println(err)
		return 0
	}
	return 2 * (len(filedata) - 1)
}

func splitName(path string) string {
	names := strings.Split(path, "\\")
	//fmt.Println(names[len(names)-1])
	// Returns the name of the file from a given path string
	return names[len(names)-1]
}

// making a function named removeDuplicate()
func removeDuplicate(arr []string) []string {
	map_var := map[string]bool{}
	result := []string{}
	for e := range arr {
		if !map_var[arr[e]] {
			map_var[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func findFile(root string, fileName string) []string {
	var files []string
	filepath.WalkDir(root, func(pathh string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if ok, err := path.Match(fileName, d.Name()); ok && err == nil {
				files = append(files, pathh)
			}
		}
		return nil
	})
	return removeDuplicate(files)
}

func compare(small_lst [][]string, big_lst [][]string, array []string) {
	var row []string
	var line []string
	var str []string
	var locations []string
	for value := range small_lst {
		row = small_lst[value]
		var name string = row[0]
		var md5 = row[1]
		var sha1 = row[2]
	Inside:
		for val := range big_lst {
			line = big_lst[val]
			if md5 == line[1] || sha1 == line[2] {
				// Now need to remove from all locations from os, but first need to ask the user's permission
				for v := range array {
					str = strings.Split(array[v], "\\")
					if str[len(str)-1] == name {
						_flag = true
						locations = findFile(global_path, name)
						fmt.Printf("\nLocations = %v\n", locations)
						// if user accept to delete
						fmt.Printf("Would you like to remove the file %v?", name)
						char, _, err := keyboard.GetSingleKey()
						if err != nil {
							panic(err)
						}
						/*
							If the answer is yes - therefore delete the file
							If the answer is no - therefore skips the file to the next one v
						*/
						res := string(char)
						switch res {
						case "y":
							fmt.Printf("\nYou pressed: y\n")
							fmt.Println(locations)
							for value := range locations {
								os.Remove(locations[value])
								fmt.Printf("file %v was removed from location %v\n", name, locations[value])
							}
							fmt.Printf("\nMoving to the next file\n")
							break Inside
						case "n":
							fmt.Printf("\nMoving to the next file\n")
							break Inside
						default:
							fmt.Println("Pressed wrong key.\nProgram terminated")
							os.Exit(1)
						}
					}
				}
			}
		}
	}
}

/*
Designed to read YAML configuration file
*Config is a pointer to the struct(Config)
*/
func readConf(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

/*
Arguments.Name - C:\Users\User\Downloads\go_testing
Arguments.Path - C:\Users\User\Downloads\go_testing\Blacklist.csv = args[2]
*/
func main() {
	conf, err := readConf("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	/*
		conf.Arguments.Path(Args[1]) used for search in path
		conf.Arguments.Name(Args[2]) used for the path and name of Blacklist.csv file
	*/
	var name string
	var path string
	flag.StringVar(&path, "path", conf.Arguments.Path, "some text description")
	flag.StringVar(&name, "name", conf.Arguments.Name, "some text description")
	args := os.Args
	fmt.Print("\n")
	fmt.Printf("All args are: %v\n", args)
	var size int = GetLengthFile((args[1])) // conf.Arguments.Name
	fmt.Printf("Total number of hashes is: %d\n", size)
	// array to repersent all rows from Blacklist.csv file with headers
	records := readCsvFile(args[2]) // conf.Arguments.Name
	var row Hash_line
	// 2d array - Each row represents a row from blacklist.csv file
	var hashes [][]string
	for value := range records {
		if records[value][0] == "Name" || records[value][1] == "MD5" || records[value][2] == "SHA1" {
			continue
		}
		row.name = records[value][0]
		row.md5 = records[value][1]
		row.sha1 = records[value][2]
		var _row = []string{}
		_row = append(_row, row.name, row.md5, row.sha1)
		hashes = append(hashes, _row)
	}
	/*
		Help array
		Contains all paths and file names from user's input of current directory to search
	*/
	array := FilePathWalkDir(args[1]) // conf.Arguments.Path
	// List which contains names and hashes of all files from given directory(recuresivly)
	var completeList [][]string
	for value := range array {
		cl := []string{}
		var fileName string = splitName(array[value])
		// fill fileName, hash_file_md5(fileName), hash_file_sha1(fileName)
		cl = append(cl, fileName, hash_file_md5(fileName), hash_file_sha1(fileName))
		completeList = append(completeList, cl)
	}

	/*
		Compare small list with bigger list
		small list is hashes
		big list is completeList
	*/
	compare(hashes, completeList, array)
	if !_flag {
		fmt.Println("\nNo files exists from Blacklist.csv")
		os.Exit(1)
	}
}
