package funsensors

import (
	"errors"

	"io/ioutil"
	"log"

	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

const w1_slave_fname = "w1_slave"

type TempReading struct {
	Id     string
	Temp_c float64
}

func ReadTemperatureFile(path string) (float64, error) {
	var temp_c float64
	var err error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return temp_c, err
	}

	lines := strings.Split(string(content), "\n")
	if strings.Contains(lines[0], "YES") && strings.Contains(lines[1], "t=") {
		i, err := strconv.ParseFloat(strings.Split(lines[1], "t=")[1], 64)
		temp_c = i / 1000.0
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = errors.New("unparseable temperature file")
	}
	return temp_c, err
}

func FindAndReadTemperatures(path string) ([]TempReading, error) {
	var err error
	out := make([]TempReading, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("error reading directory %s", path))
		return out, err
	}

	for _, file := range files {
		t_file := filepath.Join(path, file.Name(), w1_slave_fname)
		temp_c, err := ReadTemperatureFile(t_file)
		if err == nil {
			out = append(out, TempReading{file.Name(), temp_c})
		}
	}
	return out, err
}

func CentigradeToF(c float64) float64 {
	return c*1.8 + 32
}
