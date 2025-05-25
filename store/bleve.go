package store

import (
	"dojer/utils"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/blevesearch/bleve/v2"
	"github.com/mitchellh/mapstructure"
)

func getIndex() (bleve.Index, error) {
	indexPath := utils.GetDataPath("dojs.index")
	_ = utils.EnsureExists(indexPath)

	index, err := bleve.Open(indexPath)
	if err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			indexMapping := bleve.NewIndexMapping()
			index, err = bleve.New(indexPath, indexMapping)
			if err != nil {
				return nil, err
			}

			return index, nil
		}
		return nil, err
	}
	return index, nil
}

func Index(d *Doujinshi) error {
	index, err := getIndex()
	if err != nil {
		return err
	}
	defer index.Close()

	err = index.Index(d.ID, d)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromIndex(id string) error {
	index, err := getIndex()
	if err != nil {
		return err
	}
	defer index.Close()

	err = index.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

/* Just for debuging things */
func AddAll() error {
	index, err := getIndex()
	if err != nil {
		return err
	}
	defer index.Close()

	d := ListAll()
	for _, doj := range d {
		err := index.Index(doj.ID, &doj)
		if err != nil {
			return err
		}
		fmt.Printf("Id %s added to the index\n", doj.ID)
	}
	return nil
}

func BleveSearch(text string, offset int, total *int) ([]Doujinshi, error) {
	var results []Doujinshi
	index, err := getIndex()
	if err != nil {
		return results, err
	}
	// defer index.Close()
	//
	// qu := NewQuery(text).Parse()
	// qu.Dump()
	// queries := qu.ToBleveQueries()
	// query := bleve.NewConjunctionQuery(queries...)
	query := bleve.NewFuzzyQuery(text)
	search := bleve.NewSearchRequest(query)
	search.SortBy([]string{"-createdAt"})
	search.From = offset
	search.Size = LIMIT
	search.Fields = []string{"*"}
	searchResults, err := index.Search(search)

	if err != nil {
		return results, err
	}

	*total = int(searchResults.Total)
	for _, hit := range searchResults.Hits {
		var doujinshi Doujinshi
		mapstructure.Decode(hit.Fields, &doujinshi)
		results = append(results, doujinshi)
	}

	return results, nil
}

func capitalizeFields(s string) string {

	s = strings.ReplaceAll(s, ": ", ":")
	rx := regexp.MustCompile(`(\w+):`)
	result := rx.ReplaceAllStringFunc(s, func(s string) string {
		ru := []rune(s)
		ru[0] = unicode.ToUpper(ru[0])
		return string(ru)
	})

	return result
}

func preprocessQuery(query string) string {
	terms := strings.Fields(query)

	var result []string
	var quotedTerm string
	var field string
	for _, term := range terms {
		if strings.HasSuffix(term, ":") {
			field = term
			continue
		}

		if strings.HasPrefix(term, "\"") {
			quotedTerm = term
		} else if quotedTerm != "" {
			quotedTerm += " " + term

			if strings.HasSuffix(term, "\"") {
				if strings.HasPrefix(quotedTerm, "-") {
					result = append(result, field+quotedTerm)
				} else {
					result = append(result, "+"+field+quotedTerm)
				}
				quotedTerm = ""
			}
		} else if !strings.HasPrefix(term, "-") {
			result = append(result, "+"+field+term)
		} else {
			nTerm := strings.Replace(term, "-", "", 1)
			result = append(result, "-"+field+nTerm)
		}
	}

	if quotedTerm != "" {
		if strings.HasPrefix(quotedTerm, "-") {
			result = append(result, field+quotedTerm)
		} else {
			result = append(result, "+"+field+quotedTerm)
		}
	}

	return strings.Join(result, " ")
}
