package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	dis "github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	bat "github.com/Layr-Labs/eigenda/disperser/cmd/batcher/flags"
	controller "github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	enc "github.com/Layr-Labs/eigenda/disperser/cmd/encoder/flags"
	opr "github.com/Layr-Labs/eigenda/node/flags"
	churner "github.com/Layr-Labs/eigenda/operators/churner/flags"
	relay "github.com/Layr-Labs/eigenda/relay/cmd/flags"
	retriever "github.com/Layr-Labs/eigenda/retriever/flags"

	"github.com/urfave/cli"
)

var myTemplate = `
type {{.Name}} struct{
	{{range $var := .Fields}}
		{{$var.EnvVar}} string
	{{end}}
}
func (vars {{.Name}}) getEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}
 `

type ServiceConfig struct {
	Name   string
	Fields []Flag
}

type Flag struct {
	Name   string
	EnvVar string
}

func getFlag(flag cli.Flag) Flag {
	strFlag, ok := flag.(cli.StringFlag)
	if ok {
		return Flag{strFlag.Name, strFlag.EnvVar}
	}
	boolFlag, ok := flag.(cli.BoolFlag)
	if ok {
		return Flag{boolFlag.Name, boolFlag.EnvVar}
	}
	boolTFlag, ok := flag.(cli.BoolTFlag)
	if ok {
		return Flag{boolTFlag.Name, boolTFlag.EnvVar}
	}
	intFlag, ok := flag.(cli.IntFlag)
	if ok {
		return Flag{intFlag.Name, intFlag.EnvVar}
	}
	int64Flag, ok := flag.(cli.Int64Flag)
	if ok {
		return Flag{int64Flag.Name, int64Flag.EnvVar}
	}
	float64Flag, ok := flag.(cli.Float64Flag)
	if ok {
		return Flag{float64Flag.Name, float64Flag.EnvVar}
	}
	uint64Flag, ok := flag.(cli.Uint64Flag)
	if ok {
		return Flag{uint64Flag.Name, uint64Flag.EnvVar}
	}
	uintFlag, ok := flag.(cli.UintFlag)
	if ok {
		return Flag{uintFlag.Name, uintFlag.EnvVar}
	}
	durationFlag, ok := flag.(cli.DurationFlag)
	if ok {
		return Flag{durationFlag.Name, durationFlag.EnvVar}
	}
	stringSliceFlag, ok := flag.(cli.StringSliceFlag)
	if ok {
		return Flag{stringSliceFlag.Name, stringSliceFlag.EnvVar}
	}
	intSliceFlag, ok := flag.(cli.IntSliceFlag)
	if ok {
		return Flag{intSliceFlag.Name, intSliceFlag.EnvVar}
	}
	log.Fatalln("Type not found", flag)
	return Flag{}
}

func genVars(name string, flags []cli.Flag) string {

	t, err := template.New("vars").Parse(myTemplate)
	if err != nil {
		panic(err)
	}

	vars := make([]Flag, 0)
	for _, flag := range flags {
		vars = append(vars, getFlag(flag))
	}

	var doc bytes.Buffer
	err = t.Execute(&doc, ServiceConfig{name, vars})
	if err != nil {
		panic(err)
	}

	return doc.String()

}

func main() {

	configs := `package deploy

	import "reflect"
	`

	configs += genVars("DisperserVars", dis.Flags)
	configs += genVars("BatcherVars", bat.Flags)
	configs += genVars("EncoderVars", enc.Flags)
	configs += genVars("OperatorVars", opr.Flags)
	configs += genVars("RetrieverVars", retriever.Flags)
	configs += genVars("ChurnerVars", churner.Flags)
	configs += genVars("ControllerVars", controller.Flags)
	configs += genVars("RelayVars", relay.Flags)

	fmt.Println(configs)

	err := os.WriteFile("../env_vars.go", []byte(configs), 0644)
	if err != nil {
		log.Panicf("Failed to write file. Err: %s", err)
	}
}
