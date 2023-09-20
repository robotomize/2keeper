package fts

import (
	_ "embed"
)

type Document struct {
	ID   string
	Text string
	Body []byte
}

type Option func(index *Index)

func WithTokenizer(f TokenizeFunc) Option {
	return func(index *Index) {
		index.tokenizerFn = f
	}
}

func WithStemmer(f FilterFunc) Option {
	return func(index *Index) {
		index.stemmerFn = f
	}
}

func WithFilter(f FilterFunc) Option {
	return func(index *Index) {
		index.filtersFuncs = append(index.filtersFuncs, f)
	}
}

func New(opts ...Option) *Index {
	idx := &Index{
		tokenizerFn:  Tokenize,
		stemmerFn:    StemmerFilter,
		filtersFuncs: []FilterFunc{LowercaseFilter, StopwordsFilter},
		idx:          make(map[string]map[string]struct{}),
	}

	for _, o := range opts {
		o(idx)
	}

	return idx
}

type Index struct {
	tokenizerFn  TokenizeFunc
	stemmerFn    FilterFunc
	filtersFuncs []FilterFunc

	idx map[string]map[string]struct{}
}

func (i *Index) Put(documents ...Document) error {
	for _, doc := range documents {
		for _, tok := range i.analyze(doc.Text) {
			if _, ok := i.idx[tok]; !ok {
				i.idx[tok] = make(map[string]struct{})
			}

			i.idx[tok][doc.ID] = struct{}{}
		}
	}

	return nil
}

func (i *Index) Search(s string) []string {
	var docs []string
	for _, tok := range i.analyze(s) {
		if v, ok := i.idx[tok]; ok {
			tmpDocs := make([]string, 0, len(v))
			for k := range v {
				tmpDocs = append(tmpDocs, k)
			}

			if len(docs) == 0 {
				docs = append(docs, tmpDocs...)
			} else {
				docs = intersection(docs, tmpDocs)
			}
		}
	}

	return docs
}

func (i *Index) analyze(s string) []string {
	tokens := i.tokenizerFn(s)
	for _, f := range i.filtersFuncs {
		tokens = f(tokens)
	}

	return i.stemmerFn(tokens)
}

func intersection(a []string, b []string) []string {
	mx := len(a)
	if len(b) > mx {
		mx = len(b)
	}

	r := make([]string, 0, mx)
	var i, j int
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			r = append(r, a[i])
			i++
			j++
		}
	}

	return r
}
