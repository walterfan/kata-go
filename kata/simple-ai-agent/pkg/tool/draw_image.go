package tool

import (
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var baseURL = "http://www.plantuml.com/plantuml"

func GeneratePngUrl(umlText string) (string, error) {
	encoded, err := encodePlantUML(umlText)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/png/%s", baseURL, encoded)
	return url, nil
}

func DrawImage(script, imageType, outputPath string) error {

	url, err := GeneratePngUrl(script)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// encodePlantUML compresses and encodes the UML text
func encodePlantUML(text string) (string, error) {
	var b strings.Builder
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(text))
	if err != nil {
		return "", err
	}
	w.Close()

	compressed := b.String()
	return encodeToPlantUMLBase64([]byte(compressed)), nil
	//return encodeToPlantUMLBase64([]byte(text)), nil
}

// encodeToPlantUMLBase64 converts zlib-compressed bytes into PlantUML Base64 encoding
func encodeToPlantUMLBase64(data []byte) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	var encoded strings.Builder

	for i := 0; i < len(data); i += 3 {
		var b1, b2, b3 byte
		b1 = data[i]
		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		c1 := b1 >> 2
		c2 := ((b1 & 0x3) << 4) | (b2 >> 4)
		c3 := ((b2 & 0xF) << 2) | (b3 >> 6)
		c4 := b3 & 0x3F

		encoded.WriteByte(alphabet[c1])
		encoded.WriteByte(alphabet[c2])
		if i+1 < len(data) {
			encoded.WriteByte(alphabet[c3])
		}
		if i+2 < len(data) {
			encoded.WriteByte(alphabet[c4])
		}
	}
	return encoded.String()
}

func init() {

	baseURL = os.Getenv("PLANTUML_SERVER")
}
