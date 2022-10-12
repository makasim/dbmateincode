# dbmateincode

Migrate db using [dbmate](https://github.com/amacneil/dbmate) and embedded migrations.

Could be used as a workaround until dbmate fs.FS support added (see https://github.com/amacneil/dbmate/issues/193).

## Migrate

```golang
package main

import (
	"embed"
	"log"
	"net/url"

	"github.com/makasim/dbmateincode"

	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
)

//go:embed sql/*.sql
var migrationDir embed.FS

func main() {
	dbUrl, err := url.Parse("postgres://postgres:dbpass@127.0.0.1:5432/test?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	if err := dbmateincode.Migrate(dbmateincode.NewConfig(dbUrl, migrationDir)); err != nil {
		log.Fatalln(err)
	}
}
```