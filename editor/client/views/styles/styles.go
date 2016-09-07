package styles

import "github.com/davelondon/vecty"

func EditorIcon() vecty.Markup {
	return vecty.Style("color", "#cccccc")
}

func MarginRightInline() vecty.Markup {
	return vecty.Style("margin-right", "3px")
}

func MarginLeftInline() vecty.Markup {
	return vecty.Style("margin-left", "3px")
}

func DropdownIcon() vecty.Markup {
	return vecty.Style("margin-right", "5px")
}
