bt
==

bt is a very simple implementation of a behavior tree framework. It's not designed to do anything fancy. The only goal is to eliminate some boilerplate in simple use cases.

Example
-------

```
package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/DeedleFake/bt"
)

func main() {
	var name string
	ask := bt.NodeFunc(func() bt.Result {
		fmt.Printf("What's your name? ")
		fmt.Scan(&name)
		return bt.Success
	})
	greet := bt.NodeFunc(func() bt.Result {
		fmt.Printf("Nice to meet you, %v.", name)
		return bt.Success
	})
	suggest := bt.NodeFunc(func() bt.Result {
		capital := []rune(name)
		capital[0] = unicode.ToUpper(capital[0])
		fmt.Printf("Are you sure your name isn't %q?", string(capital))
		return bt.Success
	})
	notLowerDaemon := func(child bt.Node) bt.Node {
		return bt.NodeFunc(func() bt.Result {
			if strings.ToLower(name) == name {
				return bt.Failure
			}

			return child.Execute()
		})
	}

	tree := bt.Sequence(
		ask,
		bt.Selector(
			notLowerDaemon(greet),
			suggest,
		),
	)
	for bt.Run(tree) != bt.Success {
	}
}
```
