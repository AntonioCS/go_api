package Util

import (
	"github.com/gobeam/stringy"
	"log"
	"os"
	"reflect"
	"strconv"
)

func PopulateWithEnv[T any](thing T) T {

	e := reflect.ValueOf(&thing).Elem()

	//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		envName := stringy.New(varName).SnakeCase().ToUpper()
		value := os.Getenv(envName)

		if value != "" {
			switch varType.Name() {
				case "int", "int64":
					e.FieldByName(varName).SetInt(convertToInt(value))
				case "string":
					e.FieldByName(varName).SetString(value)
				default:
					panic("No suitable conversion for " + value + " of type " + varType.Name())
			}
		} else {
			log.Println("Warning: No value for ", envName)
		}
	}

	return thing
}

func convertToInt(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		panic("Unable to convert value: " + value + " to int")
	}

	return i
}