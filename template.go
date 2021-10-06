package main

type template struct {
	name, ext string
}

var templates = []template{
	{name: "Go source", ext: ".go"},
	{name: "Text file", ext: ".txt"},
	{name: "User interface", ext: ".gui.json"},
	{name: "Empty file", ext: ""},
}
