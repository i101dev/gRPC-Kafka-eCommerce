package config

import (
	"fmt"
	"os"
)

func ReadNumFromFile(filename string) (int64, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var offset int64
	_, err = fmt.Fscanf(file, "%d", &offset)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func WriteNumToFile(offset int64, filename string) error {

	// fmt.Println("\n*** >>>New offset - ", offset)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", offset)

	if err != nil {
		return err
	}

	offset += 5

	return nil
}
