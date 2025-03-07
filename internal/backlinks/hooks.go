package backlinks

import (
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/mycoopts"

	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/tools"
)

// UpdateBacklinksAfterEdit is a creation/editing hook for backlinks index
func UpdateBacklinksAfterEdit(h hyphae.Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.CanonicalName(), oldText)
	contents := fetchText(h)
	newLinks := extractHyphaLinksFromContent(h.CanonicalName(), contents)
	backlinkConveyor <- backlinkIndexEdit{h.CanonicalName(), oldLinks, newLinks, contents, oldText}
}

// UpdateBacklinksAfterDelete is a deletion hook for backlinks index
func UpdateBacklinksAfterDelete(h hyphae.Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.CanonicalName(), oldText)
	backlinkConveyor <- backlinkIndexDeletion{h.CanonicalName(), oldLinks, oldText}
}

// UpdateBacklinksAfterRename is a renaming hook for backlinks index
func UpdateBacklinksAfterRename(h hyphae.Hypha, oldName string) {
	contents := fetchText(h)
	actualLinks := extractHyphaLinksFromContent(h.CanonicalName(), contents)
	backlinkConveyor <- backlinkIndexRenaming{oldName, h.CanonicalName(), actualLinks, contents}
}

// extractHyphaLinksFromContent extracts local hypha links from the provided text.
func extractHyphaLinksFromContent(hyphaName string, contents string) []string {
	ctx, _ := mycocontext.ContextFromStringInput(contents, mycoopts.MarkupOptions(hyphaName))
	linkVisitor, getLinks := tools.LinkVisitor(ctx)
	// Ignore the result of BlockTree because we call it for linkVisitor.
	_ = mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()
	var result []string
	for _, link := range foundLinks {
		switch link := link.(type) {
		case *links.LocalLink:
			result = append(result, link.Target(ctx))
		}
	}
	return result
}
