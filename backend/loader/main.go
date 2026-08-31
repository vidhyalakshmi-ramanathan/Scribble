package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	user "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
)

func main() {
	stmts, err := gormschema.New("postgres").
		Load(
			&user.User{},
			&doc.Documents{},
			&doc.Permissions{},
			&entities.Version{},
			&user.DocumentInvites{},
			// &entities.Comments{},
		)
	if err != nil {
		panic(fmt.Sprintf("failed to load gorm schema: %v", err))
	}
	io.WriteString(os.Stdout, stmts)
}
