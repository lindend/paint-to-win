package wordlist

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"paintToWin/web"
)

type ApiWordlistInfo struct {
	Id       string
	Name     string
	Language string
	NumWords int
}

type ApiWordlist struct {
	Id       string
	Name     string
	Language string
	Words    []string
}

func GetWordlistsHandler(wordlists map[string]Wordlist) web.RequestHandler {
	wordlistInfos := make([]ApiWordlistInfo, len(wordlists))

	for _, wordlist := range wordlists {
		wordlistInfos = append(wordlistInfos, ApiWordlistInfo{
			Id:       wordlist.Id,
			Name:     wordlist.Name,
			Language: wordlist.Language,
			NumWords: len(wordlist.Words),
		})
	}

	return func(req *http.Request) (interface{}, web.ApiError) {
		return wordlistInfos, nil
	}
}

func GetWordlistHandler(wordlists map[string]Wordlist) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		vars := mux.Vars(req)
		wordlistId := vars["wordlistId"]

		wordlist, ok := wordlists[wordlistId]
		if !ok {
			return nil, web.NewApiError(http.StatusNotFound, "no such wordlist found")
		}

		return ApiWordlist{
			Id:       wordlist.Id,
			Name:     wordlist.Name,
			Language: wordlist.Language,
			Words:    wordlist.Words,
		}, nil
	}
}

func GetWordsHandler(wordlists map[string]Wordlist) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		vars := mux.Vars(req)
		wordlistId := vars["wordlistId"]
		numWordsStr, ok := vars["numWords"]
		var numWords int
		if ok {
			if nw, err := strconv.ParseInt(numWordsStr, 0, 32); err != nil {
				numWords = 20
			} else {
				numWords = int(nw)
			}
		} else {
			numWords = 20
		}

		wordlist, ok := wordlists[wordlistId]
		if !ok {
			return nil, web.NewApiError(http.StatusNotFound, "no such wordlist found")
		}

		if numWords > len(wordlist.Words) {
			return nil, web.NewApiError(http.StatusBadRequest, "more words than in the wordlist requested")
		}

		return wordlist.GetWords(numWords), nil
	}
}

func RegisterWordlistApi(router *mux.Router, wordlists []Wordlist) {
	wordlistsLookup := make(map[string]Wordlist)
	for _, wordlist := range wordlists {
		wordlistsLookup[wordlist.Id] = wordlist
	}

	router.HandleFunc("/wordlists", web.DefaultHandler(GetWordlistsHandler(wordlistsLookup))).Methods("GET", "OPTIONS")
	router.HandleFunc("/wordlists/{wordlistId}", web.DefaultHandler(GetWordlistHandler(wordlistsLookup))).Methods("GET", "OPTIONS")
	router.HandleFunc("/wordlists/{wordlistId}/words", web.DefaultHandler(GetWordsHandler(wordlistsLookup))).Methods("GET", "OPTIONS")
}
