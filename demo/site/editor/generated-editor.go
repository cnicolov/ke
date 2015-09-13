package main

import (
	_ "kego.io/demo/common/images"
	_ "kego.io/demo/common/images/types"
	_ "kego.io/demo/common/units"
	_ "kego.io/demo/common/units/types"
	_ "kego.io/demo/common/words"
	_ "kego.io/demo/common/words/types"
	_ "kego.io/demo/site"
	_ "kego.io/demo/site/types"
	"kego.io/editor/client"
	"kego.io/js/console"
	_ "kego.io/system"
	_ "kego.io/system/types"
)

func main() {
	if err := client.Start("kego.io/demo/site"); err != nil {
		console.Error(err.Error())
	}
}
