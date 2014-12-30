package wordlist

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func testMakeWordlistInfo() WordlistInfo {
	return WordlistInfo{
		Name:     "name",
		Id:       "id",
		Language: "language",
	}
}

func getNumWords(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += "word" + strconv.Itoa(i) + "\n"
	}
	return result
}

func TestLoadEmptyWordlist(t *testing.T) {
	reader := strings.NewReader("")
	info := testMakeWordlistInfo()

	wordlist := loadWordlist(info, reader)

	if info.Name != wordlist.Name {
		t.Fatal("Expected", info.Name, "found", wordlist.Name)
	}

	if info.Id != wordlist.Id {
		t.Fatal("Expected", info.Id, "found", wordlist.Id)
	}

	if info.Language != wordlist.Language {
		t.Fatal("Expected", info.Language, "found", wordlist.Language)
	}

	if len(wordlist.Words) != 0 {
		t.Fatal("Expected empty wordlist")
	}
}

func TestLoadOneWordWordlist(t *testing.T) {
	reader := strings.NewReader("word")
	info := testMakeWordlistInfo()

	wordlist := loadWordlist(info, reader)

	if len(wordlist.Words) != 1 {
		t.Fatal("Expected wordlist with one word")
	}

	if wordlist.Words[0] != "word" {
		t.Fatal("Expected word found", wordlist.Words[0])
	}
}

func TestLoadThreeWordWordlist(t *testing.T) {
	reader := strings.NewReader("word\nword2")
	info := testMakeWordlistInfo()

	wordlist := loadWordlist(info, reader)

	if len(wordlist.Words) != 2 {
		t.Fatal("Expected wordlist with two words")
	}

	if wordlist.Words[0] != "word" || wordlist.Words[1] != "word2" {
		t.Fatal("Invalid words in wordlist")
	}
}

func TestLoadWordlistWithSpace(t *testing.T) {
	reader := strings.NewReader("word with spaces")
	info := testMakeWordlistInfo()

	wordlist := loadWordlist(info, reader)

	if len(wordlist.Words) != 1 {
		t.Fatal("Expected wordlist with one word")
	}

	if wordlist.Words[0] != "word with spaces" {
		t.Fatal("Expected 'word with spaces' found", wordlist.Words[0])
	}
}

func TestLoadWordlistWithSpecialCharacters(t *testing.T) {
	word := "!\"?:;@#$%^&*(){}[]'\\/.,<>|"
	wordlist := loadWordlist(testMakeWordlistInfo(), strings.NewReader(word))

	if len(wordlist.Words) != 1 {
		t.Fatal("Expected wordlist with one word")
	}

	if wordlist.Words[0] != word {
		t.Fatal("Expected '", word, "' found '", wordlist.Words[0])
	}
}

func containsDuplicates(list []string) bool {
	items := make(map[string]*struct{})

	for _, item := range list {
		if items[item] != nil {
			return true
		}
		items[item] = &struct{}{}
	}
	return false
}

func TestGetWordlistWords(t *testing.T) {
	words := getNumWords(100)
	wordlist := loadWordlist(testMakeWordlistInfo(), strings.NewReader(words))

	for i := 0; i < 100; i++ {
		numWords := rand.Intn(30)
		chosenWords := wordlist.GetWords(numWords)

		if len(chosenWords) != numWords {
			t.Fatal("Invalid number of words returned, expected", numWords, "found", len(chosenWords), "in round", i+1)
		}

		if containsDuplicates(chosenWords) {
			t.Fatal("Returned list contained duplicates in round", i+1)
		}
	}
}

func TestGetMoreWordsThanList(t *testing.T) {
	words := getNumWords(5)
	wordlist := loadWordlist(testMakeWordlistInfo(), strings.NewReader(words))

	chosenWords := wordlist.GetWords(15)

	if len(chosenWords) != 5 {
		t.Fatal("Invalid number of words returned, expected", 5, "found", len(chosenWords))
	}

	if containsDuplicates(chosenWords) {
		t.Fatal("Returned list contained duplicates in round")
	}
}
