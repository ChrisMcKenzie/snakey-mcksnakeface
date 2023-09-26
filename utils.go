package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func PrintGameState(state GameState) {
	var buffer bytes.Buffer
	PrettyEncode(state, &buffer)
	fmt.Println(buffer.String())
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func PrettyEncode(data interface{}, out io.Writer) {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		fmt.Print("could not encode GameState")
	}

}
