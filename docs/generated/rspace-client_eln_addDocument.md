## rspace-client eln addDocument

Creates a new Basic or Structured document with optional tags and content

### Synopsis

Create a new document, with an optional name and parent folder.
If a file is a file of HTML content, it is loaded verbatim, otherwise, plain text files are wrapped in 'pre'
tags to preserve formatting.

You can also create an structured (multi-field) document by passing the 'formId'.
This creates an empty document. You can also create 1 or more documents from a CSV
file with the following characteristics:
- 1st row is a header row and is ignored
- Each row will supply data to create an RSpace document
- The column order should match the fields in the Form definition
- The total number of columns should match the total number of fields in the Form
	

```
rspace-client eln addDocument [flags]
```

### Examples

```

// create a new document with tags and HTML content
rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.html

// create a doc with tags and  plain-text content, which will be wrapped in '<pre>' tag
rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.txt

// create a doc using verbatim text
rspace eln  addDocument --name doc1  --content "some content"

// create an empty, structured (multi-field) document
rspace eln addDocument --name myDoc --formId FM2

// create a multi-field document with data in a CSV file:
rspace eln addDocument --name myDoc --formId FM2 --input data.csv


```

### Options

```
      --content string   Text or HTML content to put in a basic document
      --file string      A file of text or HTML content to put in a basic document
      --folder string    An id for the folder that will contain the new document
      --formId string    Id for a form
  -h, --help             help for addDocument
      --input string     File of input data in CSV format for adding field Data to structured documents
      --name string      A name for the document
      --tags string      One or more tags, comma separated
```

### Options inherited from parent commands

```
      --config string         config file (default is $HOME/.rspace)
  -o, --outFile string        Output file for program output
  -f, --outputFormat string   Output format: one of 'json','table', 'csv' or 'quiet'  (default "table")
```

### SEE ALSO

* [rspace-client eln](rspace-client_eln.md)	 - Top-level command to work with RSpace ELN

###### Auto generated by spf13/cobra on 14-May-2022
