package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// LineIndex returns line index of search in content.
func LineIndex(content string, search string) int {
	ret := -1
	line := ""
	for _, c := range content {
		if c == '\n' {
			line += ""
			ret++
		} else {
			line += string(c)
		}
		if strings.Index(line, search) > -1 {
			ret++
			break
		}
	}
	return ret
}

// PowerLess select titles of Power<P
func PowerLess(P int, f func(*MdTitleTree)) func(*MdTitleTree) {
	return func(n *MdTitleTree) {
		if n.Power < P && n.Title != "" {
			f(n)
		}
	}
}

// LineLess select titles of Line<P
func LineLess(P int, f func(*MdTitleTree)) func(*MdTitleTree) {
	return func(n *MdTitleTree) {
		if n.Line < P && n.Title != "" {
			f(n)
		}
	}
}

var mdTitle = regexp.MustCompile(`^([#]{1,6})\s*(.+)`)

func cntStr(in []string, what string) int {
	c := -1
	for _, i := range in {
		if i == what {
			c++
		}
	}
	if c > -1 {
		c++ // starts at 1
	}
	return c
}

// GetAllMdTitles extracts all MD titles markup.
func GetAllMdTitles(content string) []MdTitle {
	ret := []MdTitle{}
	allTitles := []string{}
	line := ""
	isInBlock := false
	isInTitle := false
	i := 0
	for _, c := range content {
		if !isInBlock && c == '\n' {
			if isInTitle {
				if mdTitle.MatchString(line) {
					got := mdTitle.FindAllStringSubmatch(line, -1)
					if len(got) > 0 {
						t := got[0][2]
						ret = append(ret, MdTitle{
							Line: i, Title: t,
							Power:     len(got[0][1]),
							Duplicate: cntStr(allTitles, t),
						})
						allTitles = append(allTitles, t)
					}
				}
			}
			i++
			isInTitle = false
			line = ""
		} else if c == '`' {
			isInBlock = !isInBlock
			line += string(c)
		} else if c == '#' && !isInBlock {
			isInTitle = true
			line += string(c)
		} else {
			if c == '\n' {
				i++
			}
			line += string(c)
		}
	}
	return ret
}

// MakeTitleTree transform a raw list of titles into a tree.
func MakeTitleTree(titles []MdTitle) *MdTitleTree {

	root := &MdTitleTree{}
	cur := root
	for _, t := range titles {

		if t.Power == 1 {
			nnew := &MdTitleTree{MdTitle: t, Parent: root}
			root.Items = append(root.Items, nnew)
			cur = nnew

		} else if t.Power > cur.Power {
			nnew := &MdTitleTree{MdTitle: t, Parent: cur}
			cur.Items = append(cur.Items, nnew)
			cur = nnew

		} else if t.Power == cur.Power {
			nnew := &MdTitleTree{MdTitle: t, Parent: cur.Parent}
			cur.Parent.Items = append(cur.Parent.Items, nnew)
			cur = nnew

		} else if t.Power < cur.Power {
			for {
				if cur.Parent.Power <= t.Power {
					break
				}
				cur = cur.Parent
			}
			cur = cur.Parent
			nnew := &MdTitleTree{MdTitle: t, Parent: cur.Parent}
			cur.Parent.Items = append(cur.Parent.Items, nnew)
			cur = nnew
		}
	}
	return root
}

// GetMdLinkHash encodes s to insert into an MD link.
func GetMdLinkHash(link string) string {
	link = strings.ToLower(link)
	link = strings.Replace(link, "/", "", -1)
	link = strings.Replace(link, "$", "", -1)
	link = strings.Replace(link, ">", "", -1)
	link = strings.Replace(link, ".", "", -1)
	link = strings.Replace(link, ";", "", -1)
	link = strings.Replace(link, ":", "", -1)
	link = strings.Replace(link, "!", "", -1)
	link = strings.Replace(link, "'", "", -1)
	link = strings.Replace(link, "|", "", -1)
	link = strings.Replace(link, "[", "", -1)
	link = strings.Replace(link, "]", "", -1)
	link = strings.Replace(link, ",", "", -1)
	link = strings.Replace(link, " ", "-", -1)
	// should it be a regexp like /[a-z0-9-_]/i ?
	return link
}

// MdTitleTree is an MdTitle with tree capabilities
type MdTitleTree struct {
	MdTitle
	Parent *MdTitleTree
	Items  []*MdTitleTree
}

// Traverse a tree
func (m *MdTitleTree) Traverse(f func(*MdTitleTree)) {
	f(m)
	for _, i := range m.Items {
		i.Traverse(f)
	}
}

func (m *MdTitleTree) String() string {
	x := strings.Repeat("#", m.Power)
	return fmt.Sprintf("%-5v %-15q Items:%v Line:%v", x, m.Title, len(m.Items), m.Line)
}

// MdTitle is a markdwon title.
type MdTitle struct {
	Line      int
	Power     int
	Duplicate int
	Title     string
}
