package store

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

var namespaces = []string{
	"id",
	"name",
	"title",
	"parodies",
	"characters",
	"tags",
	"artists",
	"groups",
	"languages",
	"categories",
	"pages",
	"uploaded",
	"createdAt",
}

type Query struct {
	OriginalText  string
	QuotedTerms   []string
	NegativeTerms []string
	PositiveTerms []string
	Fields        map[string]string
}

func (q Query) Dump() {
	fmt.Printf("OriginalText: %s\n", q.OriginalText)
	fmt.Printf("QuotedTerms : %s\n", q.QuotedTerms)
	fmt.Printf("NegativeTerm: %s\n", q.NegativeTerms)
	fmt.Printf("PositiveTerm: %s\n", q.PositiveTerms)
	fmt.Printf("Fields      : %s\n", q.Fields)
}

func NewQuery(text string) *Query {
	query := new(Query)
	query.OriginalText = text
	return query
}

func (q *Query) Parse() *Query {

	join := strings.Join(namespaces, `|`)
	pattern := regexp.MustCompile(`(` + join + `):\s?("[^"]+"|\S+)`)

	var result = make(map[string]string)

	matches := pattern.FindAllStringSubmatch(q.OriginalText, -1)

	newQuery := pattern.ReplaceAllString(q.OriginalText, "")
	for _, match := range matches {
		fieldName := match[1]
		if result[fieldName] != "" {
			result[fieldName] = result[fieldName] + "," + match[2]
		} else {
			result[match[1]] = match[2]
		}
	}

	q.Fields = result
	pattern = regexp.MustCompile(`-?(("[^"]+")|\S+)`)

	matches = pattern.FindAllStringSubmatch(newQuery, -1)

	for _, match := range matches {
		if len(match) > 2 && match[2] != "" {
			q.QuotedTerms = append(q.QuotedTerms, match[0])
			continue
		}

		if strings.HasPrefix(match[0], "-") {
			q.NegativeTerms = append(q.NegativeTerms, match[0])
		} else {
			q.PositiveTerms = append(q.PositiveTerms, match[0])
		}
	}

	return q
}

func (q *Query) ToBleveQueries() []query.Query {
	var queryStringQueryTerms []string
	var queries []query.Query

	queryStringQueryTerms = append(queryStringQueryTerms, q.NegativeTerms...)

	for _, term := range q.PositiveTerms {
		query := bleve.NewQueryStringQuery(term)
		queries = append(queries, query)
	}

	for _, term := range q.QuotedTerms {
		queryStringQueryTerms = append(queryStringQueryTerms, "+"+term)
	}

	for k, v := range q.Fields {
		multipleValues := strings.Split(v, ",")
		if len(multipleValues) > 0 {
			for _, value := range multipleValues {
				queryStringQueryTerms = append(queryStringQueryTerms, fmt.Sprintf("+%s:%s", k, value))
			}
		} else {
			queryStringQueryTerms = append(queryStringQueryTerms, fmt.Sprintf("+%s:%s", k, v))
		}
	}

	if len(queryStringQueryTerms) > 0 {
		queries = append(queries, bleve.NewQueryStringQuery(strings.Join(queryStringQueryTerms, " ")))
	}

	return queries
}
