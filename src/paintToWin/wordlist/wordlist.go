package wordlist

import (
	"bufio"
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

func loadWordlist(info WordlistInfo) (Wordlist, error) {
	file, err := os.Open(info.Path)
	if err != nil {
		return Wordlist{}, err
	}
	defer file.Close()

	wordlist := Wordlist{
		Id:       info.Id,
		Name:     info.Name,
		Language: info.Language,
		Words:    make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordlist.Words = append(wordlist.Words, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return Wordlist{}, err
	}

	return wordlist, nil
}

func enumerateWordlists(root string) ([]WordlistInfo, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	result := make([]WordlistInfo, 0)
	for _, file := range files {
		filePath := path.Join(root, file.Name())
		if file.IsDir() {
			subFiles, err := enumerateWordlists(filePath)
			if err != nil {
				result = append(result, subFiles...)
			}
		} else {
			fileParts := strings.SplitN(file.Name(), ":", 2)
			if len(fileParts) == 2 {
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
	result := make([]string, numWords)
	words := wl.Words
	for i := 0; i < numWords; i++ {
		if len(words) > 0 {
			wordIndex := rand.Intn(len(words))
			result = append(result, words[wordIndex])
			words[wordIndex] = words[0]
			words = words[1:]
		}
	}
	return result
}
