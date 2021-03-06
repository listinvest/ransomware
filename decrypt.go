package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mauri870/cryptofile/crypto"
)

func decryptFiles(key string) {

	fmt.Println("Note: \nIf you are trying a wrong key your files will be decrypted with broken content irretrievably, please don't try keys randomly\nYou have been warned")
	fmt.Println("Continue? Y/N")

	var input rune
	fmt.Scanf("%c", &input)

	if input != 'Y' {
		os.Exit(2)
	}

	log.Println("Walking dirs and searching for encrypted files...")

	// Loop over the interesting directories
	for _, f := range InterestingDirs {
		folder := BaseDir + f
		filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			ext := filepath.Ext(path)
			if ext == EncryptionExtension {
				// Matching Files encrypted
				file := File{f, ext[1:], path}
				MatchedFiles = append(MatchedFiles, file)
				log.Println("Matched:", path)
			}
			return nil
		})
	}

	// Loop over the matched files
	for _, file := range MatchedFiles {
		log.Printf("Decrypting %s...\n", file.Path)
		// Read the file content
		ciphertext, err := ioutil.ReadFile(file.Path)
		if err != nil {
			log.Println("Error opening %s\n", file.Path)
			continue
		}

		// Decrypting with the key
		text, err := crypto.Decrypt([]byte(key), ciphertext)
		if err != nil {
			log.Println(err)
			continue
		}

		// Write a new file with the decrypted content
		err = ioutil.WriteFile(file.Path[0:len(file.Path)-len(filepath.Ext(file.Path))], text, 0600)
		if err != nil {
			log.Println(err)
			continue
		}

		// Remove the encrypted file
		os.Remove(file.Path)
		if err != nil {
			log.Println(err)
		}
	}

	if len(MatchedFiles) == 0 {
		log.Println("No encrypted files found")
	}
}
