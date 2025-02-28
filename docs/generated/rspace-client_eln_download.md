## rspace-client eln download

Download attachment files by their ID

### Synopsis

Downloads 1 or more files by their ID. 'dir' flag is is optional; if not set
will download to current folder. The Ids should be for files in the Gallery. these files typically
have global ID prefix 'GL'
	

```
rspace-client eln download [flags]
```

### Examples

```

// download 3 files to current folder by their ID
rspace eln download 1234 5678 1234

// download a file to the given directory
rspace eln download 1234 --dir /downloadFolder

// globalIds work too
rspace eln download GL1234 GL12345--dir /downloadFolder
	
```

### Options

```
      --dir string   Optional directory to download into
  -h, --help         help for download
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
