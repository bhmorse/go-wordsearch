# go-wordsearch

This program takes a configuration json file containing height, width, and a list of words. It places the words, and then fills in the rest of the puzzle with random letters

The format of the json file is:

{
    "width": 10,
    "height": 10,
    "words": [
    	 "test",
	 "foo",
	 "bar",
	 "baz"
    ]
}