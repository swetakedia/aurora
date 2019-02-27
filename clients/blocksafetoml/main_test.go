package blocksafetoml

import "log"

// ExampleGetTOML gets the blocksafe.toml file for coins.asia
func ExampleClient_GetBlocksafeToml() {
	_, err := DefaultClient.GetBlocksafeToml("coins.asia")
	if err != nil {
		log.Fatal(err)
	}
}
