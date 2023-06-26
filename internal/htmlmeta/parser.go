package htmlmeta

import (
	"context"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type MetaTag struct {
	Title       string
	Description string
	Keywords    []string
}

func Parse(ctx context.Context, r io.Reader) (MetaTag, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return MetaTag{}, fmt.Errorf("html.Parse: %w", err)
	}

	var m MetaTag

	if err := traverse(ctx, doc, &m); err != nil {
		return MetaTag{}, fmt.Errorf("traverse: %w", err)
	}

	return m, nil
}

func traverse(ctx context.Context, n *html.Node, m *MetaTag) error {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "meta":
			parseMeta(m, n)
		case "title":
			if n.FirstChild != nil {
				m.Title = n.FirstChild.Data
			}
		default:
		}
	}

	done := m.Title != "" && len(m.Keywords) > 0 && m.Description != ""
	if done {
		return nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := traverse(ctx, c, m); err != nil {
			return err
		}
	}

	return nil
}

func parseMeta(m *MetaTag, n *html.Node) {
	if len(n.Attr) < 2 {
		return
	}

	var content string

	var (
		isKeywords    bool
		isDescription bool
	)

	for _, attr := range n.Attr {
		switch {
		case attr.Key == "name" && strings.ToLower(attr.Val) == "keywords":
			isKeywords = true
		case attr.Key == "name" && strings.ToLower(attr.Val) == "description":
			isDescription = true
		case attr.Key == "content":
			content = attr.Val
		}
	}

	switch {
	case isDescription:
		m.Description = content
	case isKeywords:
		tags := strings.Split(content, ",")
		for idx1 := range tags {
			tags[idx1] = strings.TrimSpace(tags[idx1])
		}

		m.Keywords = append(m.Keywords, tags...)
	default:
	}
}
