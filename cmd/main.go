package main

import (
	"fmt"

	_ "github.com/lib/pq"
	webl "github.com/zentrope/webl"
)

func main() {
	fmt.Println("Some Postgres Stuff")

	database := webl.NewDatabase("blogsvc", "wanheada", "blogdb")

	database.Connect()

	defer database.Disconnect()

	for _, a := range database.Authors() {
		fmt.Printf(" - %+v\n", a)
	}

	fmt.Println("\nauthor_types")
	for _, s := range database.AuthorTypes() {
		fmt.Printf(" - %#v\n", s)
	}

	fmt.Println("\nauthentic:")
	fmt.Printf("  - keith/test1234 -> %v\n", database.Authentic("keith", "test1234"))
	fmt.Printf("  - KEITH/test1234 -> %v\n", database.Authentic("KEITH", "test1234"))
	fmt.Printf("  - moire/dulynow  -> %v\n", database.Authentic("moire", "dylynow"))

	var keith = database.Author("keith")
	fmt.Println("\nDetail")
	fmt.Printf(" keith -> %#v\n", keith)

	fmt.Println("\nNew User")
	fmt.Printf(" mary -> exists? %#v\n", database.AuthorExists("mary"))

	database.CreateAuthor("mary", "mary@example.com", "test1234")

	fmt.Printf(" mary -> %#v\n", database.Author("mary"))
	fmt.Printf(" mary -> exists? %#v\n", database.AuthorExists("mary"))
	fmt.Printf(" mary -> deleting...\n")
	database.DeleteAuthor("mary")
	fmt.Printf(" mary -> exists? %#v\n", database.AuthorExists("mary"))

	fmt.Println("\nPosts")
	for _, p := range database.Posts() {
		fmt.Printf(" - %+v\n", p)
	}

	fmt.Println("\nKeith Posts")
	for _, p := range database.PostsByAuthor("keith") {
		fmt.Printf(" - %+v\n", p)
	}

}
