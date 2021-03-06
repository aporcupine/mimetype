package mimetype

import "mime"

// MIME struct holds information about a file format: the string representation
// of the MIME type, the extension and the parent file format.
//
// This struct should not be instantiated by clients. It is exposed purely for
// documentation purposes.
type MIME struct {
	mime      string
	aliases   []string
	extension string
	matchFunc func([]byte) bool
	children  []*MIME
	parent    *MIME
}

// String returns the string representation of the MIME type, e.g., "application/zip".
func (n *MIME) String() string {
	return n.mime
}

// Extension returns the file extension associated with the MIME type.
// It includes the leading dot, as in ".html". When the file format does not
// have an extension, the empty string is returned.
func (n *MIME) Extension() string {
	return n.extension
}

// Parent returns the parent MIME type from the hierarchy.
// Each MIME type has a non-nil parent, except for the root MIME type.
//
// For example, the application/json and text/html MIME types have text/plain as
// their parent because they are text files who happen to contain JSON or HTML.
// Another example is the ZIP format, which is used as container
// for Microsoft Office files, EPUB files, JAR files and others.
func (n *MIME) Parent() *MIME {
	return n.parent
}

// Is checks whether this MIME type, or any of its aliases, is equal to the
// expected MIME type. MIME type equality test is done on the "type/subtype"
// sections, ignores any optional MIME parameters, ignores any leading and
// trailing whitespace, and is case insensitive.
func (n *MIME) Is(expectedMIME string) bool {
	// Parsing is needed because some detected MIME types contain parameters
	// that need to be stripped for the comparison.
	expectedMIME, _, _ = mime.ParseMediaType(expectedMIME)
	found, _, _ := mime.ParseMediaType(n.mime)

	if expectedMIME == found {
		return true
	}
	for _, alias := range n.aliases {
		if alias == expectedMIME {
			return true
		}
	}

	return false
}

func newMIME(mime, extension string, matchFunc func([]byte) bool, children ...*MIME) *MIME {
	n := &MIME{
		mime:      mime,
		extension: extension,
		matchFunc: matchFunc,
		children:  children,
	}

	for _, c := range children {
		c.parent = n
	}

	return n
}

func (n *MIME) alias(aliases ...string) *MIME {
	n.aliases = aliases
	return n
}

// match does a depth-first search on the matchers tree.
// It returns the deepest successful matcher for which all the children fail.
func (n *MIME) match(in []byte, deepestMatch *MIME) *MIME {
	for _, c := range n.children {
		if c.matchFunc(in) {
			return c.match(in, c)
		}
	}

	return deepestMatch
}

func (n *MIME) flatten() []*MIME {
	out := []*MIME{n}
	for _, c := range n.children {
		out = append(out, c.flatten()...)
	}

	return out
}
