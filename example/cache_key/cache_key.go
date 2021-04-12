package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/latavin243/tagcollector"
)

// SearchBookRequest defines the request when some user search for a book
// Cache for these requests are needed
// Cache key contains UserID, Author and BookName in order
type SearchBookRequest struct {
	UserID           string `cache:"1"`
	BookName         string `cache:"3"`
	Author           string `cache:"2"`
	RequestTimeStamp int64
	RequestID        string
}

func (r *SearchBookRequest) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *r)
}

type CacheKeySegment struct {
	Order int
	Key   string
}

var (
	sampleReq = &SearchBookRequest{
		UserID:           "123123",
		BookName:         "The Three-Body Problem",
		Author:           "Ken Liu",
		RequestTimeStamp: 1618241420,
		RequestID:        "1fa6b2bd6c154c72fab82762d6951809",
	}

	cacheKeySeparator = ":"
	cacheKeyTag       = "cache"
)

func getCacheKey(req interface{}, separator string) (cacheKey string, err error) {
	cacheKeyEntries, err := tagcollector.Collect(req, []string{cacheKeyTag})
	if err != nil {
		return "", err
	}

	cacheKeySegments := make([]*CacheKeySegment, 0)
	for _, entry := range cacheKeyEntries {
		rawOrderStr, ok := entry.TagMap[cacheKeyTag]
		if !ok {
			continue
		}
		order, err := strconv.Atoi(rawOrderStr)
		if err != nil {
			return "", fmt.Errorf("convert cache key order string to int error, rawCacheKeyOrder=%s, err=%s",
				rawOrderStr, err)
		}
		cacheKeySegments = append(cacheKeySegments, &CacheKeySegment{
			Order: order,
			Key:   fmt.Sprintf("%+v", entry.FieldValue),
		})
	}

	sort.Slice(cacheKeySegments, func(i, j int) bool {
		return cacheKeySegments[i].Order < cacheKeySegments[j].Order
	})

	cacheKeyStrSegments := make([]string, 0)
	for _, cacheKeySegment := range cacheKeySegments {
		cacheKeyStrSegments = append(cacheKeyStrSegments, cacheKeySegment.Key)
	}

	return strings.Join(cacheKeyStrSegments, cacheKeySeparator), nil
}

func main() {
	fmt.Printf("raw request is: %s\n", sampleReq)
	cacheKey, err := getCacheKey(sampleReq, cacheKeySeparator)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cache key is: `%s`\n", cacheKey)
}
