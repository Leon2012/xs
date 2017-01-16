package xs

import (
	_ "fmt"
	"os"
	_ "strconv"
	"strings"
)

func LoadInIFile(filePath string) (map[string]interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fi.Size()
	data := make([]byte, fileSize)
	_, err = f.Read(data)
	if err != nil {
		return nil, err
	}
	ini := ParseIniData(string(data))
	return ini, nil
}

func ParseIniData(content string) map[string]interface{} {
	var ret map[string]interface{}
	ret = make(map[string]interface{})
	lines := strings.Split(content, "\n")
	var cur map[string]string
	var prevSec string
	prevSec = ""
	for _, line := range lines {
		if line == "" || line[0:1] == ";" || line[0:1] == "#" {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sec := line[1:(len(line) - 1)]
			_, ok := ret[sec]
			if !ok {
				cur = make(map[string]string)
				ret[sec] = cur
			}
			prevSec = sec
			continue
		}
		pos := strings.Index(line, "=")
		if pos == -1 {
			continue
		}
		key := strings.TrimSpace(line[0:pos])
		value := line[(pos + 1):(len(line))]
		value = strings.Trim(value, " '\t\"")
		//fmt.Println("key:" + key + " value:" + value + " pos:" + strconv.Itoa(pos) + " len:" + strconv.Itoa(len(line)))
		//fmt.Println(prevSec)
		if prevSec != "" {
			m, ok := ret[prevSec]
			if ok {
				if cur, ok = m.(map[string]string); ok {
					cur[key] = value
					ret[prevSec] = cur
				}
			}
		} else {
			ret[key] = value
			// apos := strings.Index(key, ".")
			// aSec := key[0:apos]
			// cur, ok := ret[aSec]
			// if !ok {
			// 	cur = make(map[string]string)
			// }
			// aKey := key[(apos + 1):len(key)]
			// aValue := value
			// cur[aKey] = aValue
			// ret[aSec] = cur
		}
	}
	return ret
}
