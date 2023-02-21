package trace

import "github.com/oklog/ulid/v2"

func Generate() string {
	return ulid.Make().String()
}
