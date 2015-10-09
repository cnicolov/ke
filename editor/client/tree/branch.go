package tree // import "kego.io/editor/client/tree"

import (
	"time"

	"honnef.co/go/js/dom"
	"kego.io/kerr"
	"kego.io/system"
)

type Branch struct {
	Tree     *Tree
	root     bool
	item     Item
	parent   *Branch
	children []*Branch
	index    int
	siblings []*Branch
	next     *Branch
	prev     *Branch
	open     bool
	loading  bool
	level    int
	element  *dom.HTMLDivElement
	opener   *dom.HTMLAnchorElement
	inner    *dom.HTMLDivElement
	content  *dom.HTMLDivElement
	selected bool
	editor   system.Editor
}

func (b *Branch) Initialise() {
	// We must tolerate having a nil dom element in order to run tests in pure go
	if b.element != nil {
		b.item.Initialise(b.content)
	}
}

func (b *Branch) Append(child *Branch) *Branch {
	child.parent = b
	// We must tolerate passing in a nil dom element in order to run tests in pure go
	if b.element != nil {

		opener := dom.GetWindow().Document().CreateElement("a").(*dom.HTMLAnchorElement)
		opener.SetAttribute("class", "toggle")
		opener.AddEventListener("click", true, func(e dom.Event) {
			if child.canOpen() {
				child.Toggle()
			} else {
				child.Select(false)
			}
		})

		contentDiv := dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)
		contentDiv.SetAttribute("class", "content")

		contentDiv.AddEventListener("click", true, func(e dom.Event) {
			child.Select(false)
		})

		innerDiv := dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)
		innerDiv.SetAttribute("class", "children")
		// children should be hidden by default
		innerDiv.Style().Set("display", "none")

		nodeDiv := dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)
		nodeDiv.SetAttribute("class", "node")

		nodeDiv.AppendChild(opener)
		nodeDiv.AppendChild(contentDiv)
		nodeDiv.AppendChild(innerDiv)

		child.opener = opener
		child.element = nodeDiv
		child.inner = innerDiv
		child.content = contentDiv

		b.inner.AppendChild(nodeDiv)
	}
	b.children = append(b.children, child)
	child.Initialise()
	b.update()
	return child
}

func (b *Branch) Select(fromKeyboard bool) {

	if b.Tree.selected != nil {
		b.Tree.selected.Unselect()
	}
	b.content.Class().Add("selected")
	b.Tree.selected = b
	b.selected = true

	if fromKeyboard && b.isAsyncAndNotLoaded() {
		// wait 50ms before showing the edit panel so we don't generate
		// lots of content load requests if we scroll quickly. We only
		// do this for keyboard events.
		go func() {
			time.Sleep(time.Millisecond * 50)
			if b.selected {
				b.showEditPanel(fromKeyboard)
			}
		}()
	} else {
		b.showEditPanel(fromKeyboard)
	}

	return
}

func (b *Branch) showEditPanel(fromKeyboard bool) {

	if b.root {
		return
	}

	if b.editor != nil && b.Tree.editor != nil && b.editor == b.Tree.editor {
		return
	}

	if b.Tree.editor != nil {
		b.Tree.editor.Hide()
		b.Tree.editor = nil
	}

	done, ok := b.ensureContentLoaded()

	if !ok {
		// if the operation should be cancelled, we should return immediately
		return
	}

	success := func() {

		if !fromKeyboard && b.canOpen() && !b.open {
			// if we clicked on an item, and it's not open, we should open it
			b.Open()
		}

		if b.editor == nil {
			hn, ok := b.item.(HasNode)
			if !ok {
				return
			}
			n := hn.Node()
			he, ok := n.Value.(system.HasEditor)
			if ok {
				b.editor = he.GetEditor(n)
			} else {
				b.editor = n.GetEditor()
			}
		}
		if b.editor == nil {
			return
		}
		if !b.editor.Initialized() {
			panel := dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)
			panel.Style().Set("display", "none")
			panel.Class().SetString("mdl-color--white mdl-shadow--2dp mdl-cell mdl-cell--12-col mdl-grid")
			b.Tree.content.AppendChild(panel)
			if err := b.editor.Initialize(panel, b.Tree.Path, b.Tree.Aliases); err != nil {
				b.Tree.Fail <- kerr.New("KKOBKWJDBI", err, "Initialize")
				return
			}

		}
		if b.Tree.editor != nil {
			b.Tree.editor.Hide()
		}
		b.editor.Show()
		b.Tree.editor = b.editor

	}

	if done == nil {
		// if the done chanel is nil, the operation was synchronous, so we should call success synchronously
		success()
	} else {
		go func() {
			// block and wait until the response arrives
			<-done
			success()
		}()
	}
	return
}

func (b *Branch) Unselect() {
	// un-select
	b.content.Class().Remove("selected")
	b.Tree.selected = nil
	b.selected = false
}

func (b *Branch) Each(f func(*Branch) error) error {
	c := b
	for c != nil {
		if err := f(c); err != nil {
			return kerr.New("MPQEUJHVPN", err, "f")
		}
		c = c.next
	}
	return nil
}

func (b *Branch) ensureContentLoaded() (done chan bool, success bool) {

	done = make(chan bool, 1)

	if async, ok := b.item.(AsyncItem); ok && !async.ContentLoaded() {

		if b.loading {
			// if we're already in the process of loading the contents, we should
			// cancel the operation
			return nil, false
		}

		// load content asynchronously
		b.loading = true
		responseChannel := async.LoadContent()

		go func() {
			<-responseChannel
			b.loading = false
			done <- true
		}()

		return done, true

	} else {
		// if item is not async or content is already loaded, just
		// open the node.
		return nil, true
	}
}

func (b *Branch) Open() {

	if b.root {
		return
	}

	done, ok := b.ensureContentLoaded()

	if !ok {
		// if the operation should be cancelled, we should return immediately
		return
	}

	success := func() {
		if b.inner != nil {
			b.inner.Style().Set("display", "block")
		}
		b.open = true
		b.afterStateChange()
	}

	if done == nil {
		// if the done chanel is nil, the operation was synchronous, so we should call success synchronously
		success()
	} else {
		go func() {
			// block and wait until the response arrives
			<-done
			success()
		}()
	}
	return
}

func (b *Branch) Close() {
	if b.root {
		return
	}
	if b.inner != nil {
		b.inner.Style().Set("display", "none")
	}
	b.open = false
	b.loading = false
	b.afterStateChange()
	return
}

func (b *Branch) Toggle() {
	if b.open {
		b.Close()
	} else {
		b.Open()
	}
}

func (b *Branch) afterStateChange() {
	b.update()
	if next := b.nextVisible(); next != nil {
		// we must also update the next in the list to ensure it's prev
		// property is set correctly
		next.update()
	}
	if b.Tree.selected != nil && !b.Tree.selected.IsVisible() {
		// if the selected branch is now invisible, we should un-select it.
		b.Tree.selected.Unselect()
	}
}

// update assumes parent and index are sources of truth, and updates
// siblings, prev and next, and updates the children
func (b *Branch) update() {
	if b.parent == nil {
		// special case for the root node
		b.siblings = []*Branch{b}
		b.level = 0
	} else {
		b.siblings = b.parent.children
		b.level = b.parent.level + 1
	}
	if b.index == 0 {
		// at the start of a list of children, the previous in list
		// order is always the parent
		b.prev = b.parent
	} else {
		// in the middle of a list of children, the previous in list
		// order is the lastVisible of the previous sibling.
		b.prev = b.parent.children[b.index-1].lastVisible()
	}
	if b.open && len(b.children) > 0 {
		// if the node is open, the next in list order will be the
		// first child
		b.next = b.children[0]
	} else {
		// if it's closed, we use the nextVisible method to get the
		// next sibling, or if we're the last sibling, the next sibling
		// of the first ancestor that has one.
		b.next = b.nextVisible()
	}
	if b.open && len(b.children) > 0 {
		for index, child := range b.children {
			child.parent = b
			child.index = index
			child.update()
		}
	}
	if b.opener != nil {

		plus := `<svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M0 0h24v24H0z" fill="none"/><path d="M13 7h-2v4H7v2h4v4h2v-4h4v-2h-4V7zm-1-5C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/></svg>`
		minus := `<svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M0 0h24v24H0z" fill="none"/><path d="M7 11v2h10v-2H7zm5-9C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/></svg>`
		point := `<svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z"/><path d="M0 0h24v24H0z" fill="none"/></svg>`

		if b.isAsyncAndNotLoaded() {
			b.opener.SetInnerHTML(plus)
		} else if b.children == nil || len(b.children) == 0 {
			b.opener.SetInnerHTML(point)
		} else if b.open {
			b.opener.SetInnerHTML(minus)
		} else {
			b.opener.SetInnerHTML(plus)
		}
	}
}

func (b *Branch) isAsyncAndNotLoaded() bool {
	async, isAsync := b.item.(AsyncItem)
	if isAsync && !async.ContentLoaded() {
		return true
	}
	return false
}

func (b *Branch) canOpen() bool {
	if b.isAsyncAndNotLoaded() {
		return true
	}
	if len(b.children) == 0 {
		return false
	}
	return true
}

// lastVisible returns the last visible descendant in list order
func (b *Branch) lastVisible() *Branch {
	i := b
	for i.open && len(i.children) > 0 {
		// if the node is open, test it's last child
		i = i.children[len(i.children)-1]
	}
	// return the first node we find that's closed
	return i
}

// nextVisible returns the next sibling. If we're the last sibling,
// we find the nearest ancestor that has a next sibling, or nil if
// we're at the end of the tree.
func (b *Branch) nextVisible() *Branch {
	i := b
	for i.index >= len(i.siblings)-1 {
		// if the node is the last of the siblings,
		// test it's parent
		if i.root {
			// if we get to the root node, we're testing the last visible
			// node, so we return nil
			return nil
		}
		i = i.parent
	}
	// return the next sibling of the first ancestor that has one
	return i.siblings[i.index+1]
}

func (b *Branch) IsDescendantOf(ancestor *Branch) bool {
	test := b.parent
	for test != nil {
		if test == ancestor {
			return true
		}
		test = test.parent
	}
	return false
}

func (b *Branch) IsAncestorOf(child *Branch) bool {
	test := child.parent
	for test != nil {
		if test == b {
			return true
		}
		test = test.parent
	}
	return false
}

func (b *Branch) IsVisible() bool {
	test := b.parent
	for test != nil {
		if !test.open {
			return false
		}
		test = test.parent
	}
	return true
}
