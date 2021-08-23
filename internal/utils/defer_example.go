package utils

import (
	"fmt"
	"io"
	"os"
)

func OpenFileInCycle(path string){
	for i := 0; i < 10; i++ {
		content, err := func() (string, error) {
			f, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return "", err
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("Error when closing file:", err)
				}
			}(f)

			currContent := ""
			tmp := make([]byte, 100)
			for {
				_, err := f.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("Error while reading file:", err)
						return currContent, err
					}
					break
				}
				currContent += string(tmp)
			}

			return currContent, nil
		}()

		if err != nil {
			fmt.Printf("Error when opening file: %s", err)
			return
		}
		fmt.Printf("File content: `%s`\n", content)
	}
}
