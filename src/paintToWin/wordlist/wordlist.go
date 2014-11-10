package wordlist

import (
	"bufio"
	"ioutil"
)

type Wordlist struct {
	Name     string
	Language string
	Words    []string
}

type WordlistInfo struct {
	Name     string
	Language string
}

func loadWordlist(path string) (Wordlist, error) {
	return
}

func enumerateWordLists(root string) ([]string, error) {

}
