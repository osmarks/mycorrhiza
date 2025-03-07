// Package backlinks maintains the index of backlinks and lets you update it and query it.
package backlinks

import (
	"log/slog"
	"os"
	"sort"
	"time"
	"strings"
	"unicode"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/rivo/uniseg"
	"github.com/dchest/stemmer/porter2"
)

// yieldHyphaBacklinks gets backlinks for the desired hypha, sorts and yields them one by one.
func yieldHyphaBacklinks(hyphaName string) <-chan string {
	hyphaName = util.CanonicalName(hyphaName)
	out := make(chan string)
	sorted := hyphae.PathographicSort(out)
	go func() {
		backlinks, exists := backlinkIndex[hyphaName]
		if exists {
			for link := range backlinks {
				out <- link
			}
		}
		close(out)
	}()
	return sorted
}

var backlinkConveyor = make(chan backlinkIndexOperation) // No need to buffer because these operations are rare.

// RunBacklinksConveyor runs an index operation processing loop. Call it somewhere in main.
func RunBacklinksConveyor() {
	// It is supposed to run as a goroutine for all the time. So, don't blame the infinite loop.
	defer close(backlinkConveyor)
	for {
		(<-backlinkConveyor).apply()
	}
}

var ZeroTime time.Time

type Metadata struct {
	Outlinks []string
	Bytes    int
	Words    int
	Updated  time.Time
}

var backlinkIndex = make(map[string]linkSet)
var forwardIndex = make(map[string]Metadata)
var invertedIndex = make(map[string]map[string]int)

func scrubInvertedIndexEntry(hyphaName string, tokens []string) {
	for _, token := range tokens {
		if tmap, exists := invertedIndex[token]; exists {
			delete(tmap, hyphaName)
		}
	}
}

func writeTokensToInvertedIndex(hyphaName string, tokens []string) {
	for _, token := range tokens {
		if tmap, exists := invertedIndex[token]; exists {
			tmap[hyphaName] += 1
		} else {
			tmap := make(map[string]int)
			tmap[hyphaName] = 1
			invertedIndex[token] = tmap
		}
	}
}

func containsAlnum(s string) bool {
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			return true
		}
	}

    return false
}

func tokenize(s string) []string {
	eng := porter2.Stemmer

	tokens := make([]string, 0)
	state := -1
	remainder := s
	var c string
	for len(remainder) > 0 {
		c, remainder, state = uniseg.FirstWordInString(remainder, state)
		if containsAlnum(c) {
			token := eng.Stem(strings.ToLower(c))
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func generateMetadata(h hyphae.Hypha, content string) {
	tokens := tokenize(content)

	meta := Metadata{
		Outlinks: extractHyphaLinksFromContent(h.CanonicalName(), content),
		Bytes:    len(content),
		Words:    len(tokens),
	}
	forwardIndex[h.CanonicalName()] = meta
}

func updateRevTimestamp(h hyphae.Hypha, newTime time.Time) {
	if _, exists := forwardIndex[h.CanonicalName()]; exists {
		// ??? Golang ?????
		meta := forwardIndex[h.CanonicalName()]
		meta.Updated = newTime
		forwardIndex[h.CanonicalName()] = meta
	}
}

// IndexBacklinks traverses all text hyphae, extracts links from them and forms an initial index. Call it when indexing and reindexing hyphae.
func IndexBacklinks() {
	// It is safe to ignore the mutex, because there is only one worker.
	for h := range hyphae.YieldExistingHyphae() {
		content := fetchText(h)
		foundLinks := extractHyphaLinksFromContent(h.CanonicalName(), content)
		for _, link := range foundLinks {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][h.CanonicalName()] = struct{}{}
		}
		generateMetadata(h, content)
		if revs, err := history.Revisions(h.CanonicalName()); err == nil {
			// sorted newest first
			if len(revs) > 0 {
				updateRevTimestamp(h, revs[0].Time)
			}
		}

		writeTokensToInvertedIndex(h.CanonicalName(), tokenize(util.BeautifulName(h.CanonicalName())))
		writeTokensToInvertedIndex(h.CanonicalName(), tokenize(content))
	}
}

func Search(query string) []string {
	tokens := tokenize(query)
	result := make(map[string]int)
	for _, token := range tokens {
		if documents, exists := invertedIndex[token]; exists {
			for name, termFrequency := range documents {
				result[name] += termFrequency
			}
		}
	}
	// TODO: actually use the tf
	sortedResult := make([]string, 0)
	for name, _ := range result {
		sortedResult = append(sortedResult, name)
	}
	sort.Strings(sortedResult)
	return sortedResult
}

// BacklinksCount returns the amount of backlinks to the hypha. Pass canonical names.
func BacklinksCount(hyphaName string) int {
	if links, exists := backlinkIndex[hyphaName]; exists {
		return len(links)
	}
	return 0
}

func BacklinksFor(hyphaName string) []string {
	var backlinks []string
	for b := range yieldHyphaBacklinks(hyphaName) {
		backlinks = append(backlinks, b)
	}
	return backlinks
}

func MetadataFor(hyphaName string) Metadata {
	return forwardIndex[hyphaName]
}

func Orphans() []string {
	var orphans []string
	for h := range hyphae.YieldExistingHyphae() {
		if BacklinksCount(h.CanonicalName()) == 0 {
			orphans = append(orphans, h.CanonicalName())
		}
	}
	sort.Strings(orphans)
	return orphans
}

// Using set here seems like the most appropriate solution
type linkSet map[string]struct{}

func toLinkSet(xs []string) linkSet {
	result := make(linkSet)
	for _, x := range xs {
		result[x] = struct{}{}
	}
	return result
}

func fetchText(h hyphae.Hypha) string {
	var path string
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		return ""
	case *hyphae.TextualHypha:
		path = h.TextFilePath()
	case *hyphae.MediaHypha:
		if !h.HasTextFile() {
			return ""
		}
		path = h.TextFilePath()
	}

	text, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read file", "path", path, "err", err, "hyphaName", h.CanonicalName())
		return ""
	}
	return string(text)
}

// backlinkIndexOperation is an operation for the backlink index. This operation is executed async-safe.
type backlinkIndexOperation interface {
	apply()
}

// backlinkIndexEdit contains data for backlink index update after a hypha edit
type backlinkIndexEdit struct {
	name     string
	oldLinks []string
	newLinks []string
	content string
	oldContent string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexEdit) apply() {
	oldLinks := toLinkSet(op.oldLinks)
	newLinks := toLinkSet(op.newLinks)
	for link := range oldLinks {
		if _, exists := newLinks[link]; !exists {
			delete(backlinkIndex[link], op.name)
		}
	}
	for link := range newLinks {
		if _, exists := oldLinks[link]; !exists {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][op.name] = struct{}{}
		}
	}
	hyp := hyphae.ByName(op.name)
	generateMetadata(hyp, op.content)
	// wrong, but close enough
	updateRevTimestamp(hyp, time.Now())

	scrubInvertedIndexEntry(op.name, tokenize(op.oldContent))
	writeTokensToInvertedIndex(op.name, tokenize(op.content))
}

// backlinkIndexDeletion contains data for backlink index update after a hypha deletion
type backlinkIndexDeletion struct {
	name   string
	links  []string
	content string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexDeletion) apply() {
	for _, link := range op.links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.name)
		}
	}
	delete(forwardIndex, op.name)

	scrubInvertedIndexEntry(op.name, tokenize(op.content))
	scrubInvertedIndexEntry(op.name, tokenize(util.BeautifulName(op.name)))
}

// backlinkIndexRenaming contains data for backlink index update after a hypha renaming
type backlinkIndexRenaming struct {
	oldName string
	newName string
	links   []string
	content string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexRenaming) apply() {
	for _, link := range op.links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.oldName)
			backlinkIndex[link][op.newName] = struct{}{}
		}
	}

	scrubInvertedIndexEntry(op.oldName, tokenize(op.content))
	scrubInvertedIndexEntry(op.oldName, tokenize(util.BeautifulName(op.oldName)))
	writeTokensToInvertedIndex(op.newName, tokenize(op.content))
	writeTokensToInvertedIndex(op.newName, tokenize(util.BeautifulName(op.newName)))
}
