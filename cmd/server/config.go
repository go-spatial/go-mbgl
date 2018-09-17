package main

type RasterConfig struct {
	Hostname string `toml:"hostname"`
	Maps     []Map
}

type Map struct {
	Name  string `toml:"name"`
	Style string `style:"style"`
}
