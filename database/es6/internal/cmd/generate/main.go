// Licensed to Elasticsearch B.V. under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package main

import (
	"jinycoo.com/jinygo/database/es6/internal/cmd/generate/commands"
	_ "jinycoo.com/jinygo/database/es6/internal/cmd/generate/commands/gensource"
	_ "jinycoo.com/jinygo/database/es6/internal/cmd/generate/commands/genstruct"
	_ "jinycoo.com/jinygo/database/es6/internal/cmd/generate/commands/gentests"
)

func main() {
	commands.Execute()
}
