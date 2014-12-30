package api

import (
	"errors"
	"fmt"

	"paintToWin/service"
	"paintToWin/wordlist"
)

const Service = "wordlist"

type ApiWordlistInfo struct {
	Id       string
	Name     string
	Language string
	NumWords int
}

type GetWordlistsOutput struct {
	Wordlists []ApiWordlistInfo
}

type GetWordlistInput struct {
	WordlistId string
}

type GetWordlistOutput struct {
	Id       string
	Name     string
	Language string
	Words    []string
}

type GetWordsInput struct {
	WordlistId string
	NumWords   int
}

type GetWordsOutput struct {
	Words []string
}

var GetWordlistsOperation = service.NewOperation(Service, "getWordlists", "wordlists", "GET", nil, GetWordlistsOutput{})
var GetWordlistOperation = service.NewOperation(Service, "getWordlist", "wordlists/{WordlistId}", "GET", GetWordlistInput{}, GetWordlistOutput{})
var GetWordsOperation = service.NewOperation(Service, "getWords", "wordlists/{WordlistId}/words", "GET", GetWordsInput{}, GetWordsOutput{})

func registerGetWordlists(host service.Host, wordlists map[string]wordlist.Wordlist) {
	wordlistInfos := make([]ApiWordlistInfo, 0)

	for _, wordlist := range wordlists {
		wordlistInfos = append(wordlistInfos, ApiWordlistInfo{
			Id:       wordlist.Id,
			Name:     wordlist.Name,
			Language: wordlist.Language,
			NumWords: len(wordlist.Words),
		})
	}

	host.Register(func() GetWordlistsOutput {
		fmt.Println("Returning word lists:", wordlistInfos)
		return GetWordlistsOutput{
			Wordlists: wordlistInfos,
		}
	}, GetWordlistsOperation)
}

func registerGetWordlist(host service.Host, wordlists map[string]wordlist.Wordlist) {
	host.Register(func(input GetWordlistInput) (GetWordlistOutput, error) {
		wordlist, ok := wordlists[input.WordlistId]
		if !ok {
			return GetWordlistOutput{}, errors.New("No such wordlist found")
		}

		return GetWordlistOutput{
			Id:       wordlist.Id,
			Name:     wordlist.Name,
			Language: wordlist.Language,
			Words:    wordlist.Words,
		}, nil
	}, GetWordlistOperation)
}

func registerGetWords(host service.Host, wordlists map[string]wordlist.Wordlist) {
	host.Register(func(input GetWordsInput) (GetWordsOutput, error) {
		numWords := input.NumWords

		if numWords <= 0 {
			numWords = 20
		}

		wordlist, ok := wordlists[input.WordlistId]
		if !ok {
			return GetWordsOutput{}, errors.New("No such wordlist found")
		}

		if numWords > len(wordlist.Words) {
			return GetWordsOutput{}, errors.New("Requested more words from list than available")
		}

		fmt.Println("Returning", numWords, "from", input.WordlistId)

		return GetWordsOutput{
			Words: wordlist.GetWords(numWords),
		}, nil
	}, GetWordsOperation)
}

func RegisterWordlistApi(host service.Host, wordlists []wordlist.Wordlist) {
	wordlistsLookup := make(map[string]wordlist.Wordlist)
	for _, wordlist := range wordlists {
		wordlistsLookup[wordlist.Id] = wordlist
	}

	registerGetWordlists(host, wordlistsLookup)
	registerGetWordlist(host, wordlistsLookup)
	registerGetWords(host, wordlistsLookup)
}
