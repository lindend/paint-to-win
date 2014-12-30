package wordlist

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
)

type Wordlist struct {
	Id       string
	Name     string
	Language string
	Words    []string
}

type WordlistInfo struct {
	Id       string
	Name     string
	Language string
	Path     string
}

func LoadWordlistFromFile(info WordlistInfo) (Wordlist, error) {
	file, err := os.Open(info.Path)
	if err != nil {
		return Wordlist{}, err
	}
	defer file.Close()

	return loadWordlist(info, file), nil
}

func loadWordlist(info WordlistInfo, file io.Reader) Wordlist {
	wordlist := Wordlist{
		Id:       info.Id,
		Name:     info.Name,
		Language: info.Language,
	}

	wordlist.Words = loadWords(file)

	return wordlist
}

func loadWords(reader io.Reader) []string {
	result := make([]string, 0)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}

func EnumerateWordlists(root string) ([]WordlistInfo, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	result := make([]WordlistInfo, 0)
	for _, file := range files {
		filePath := path.Join(root, file.Name())
		if file.IsDir() {
			if file.Name()[0] != '.' {
				subFiles, err := EnumerateWordlists(filePath)
				if err == nil {
					result = append(result, subFiles...)
				}
			}
		} else {
			fileParts := strings.SplitN(file.Name(), ".", 2)
			if len(fileParts) == 2 {
				fmt.Println("Found wordlist!", file.Name())
				result = append(result, WordlistInfo{
					Id:       file.Name(),
					Name:     fileParts[0],
					Language: fileParts[1],
					Path:     path.Join(filePath),
				})
			}
		}
	}

	return result, nil
}

func (wl Wordlist) GetWords(numWords int) []string {
	result := make([]string, 0)
	selectedWords := make(map[string]struct{})

	if numWords >= len(wl.Words) {
		return wl.Words
	}

	for len(result) < numWords {
		word := wl.Words[rand.Intn(len(wl.Words))]
		if _, ok := selectedWords[word]; !ok {
			selectedWords[word] = struct{}{}
			result = append(result, word)
		}
	}
	return result
}
