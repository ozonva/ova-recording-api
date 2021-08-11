package utils

import (
	"fmt"
	"os"
)

func OpenFileInCycle(path string) {
	for i := 0; i < 10; i++ {
		content := func() string {
			f, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return ""
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("Error when closing file:", err)
				}
			}(f)

			tmp := make([]byte, 100)
			numRead, err := f.Read(tmp)
			if err != nil {
				fmt.Println("Error while reading file:", err)
				return ""
			}

			return string(tmp[:numRead])
		}()
		fmt.Printf("File content: `%s`\n", content)
	}
}
