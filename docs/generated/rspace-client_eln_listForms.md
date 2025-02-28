## rspace-client eln listForms

Lists forms

### Synopsis

List forms, sorted or paginated

```
rspace-client eln listForms [flags]
```

### Examples

```

rspace eln listForms --orderBy name --maxResults 100
	
```

### Options

```
  -h, --help               help for listForms
      --maxResults int     Maximum number of results to retrieve (default 20)
      --orderBy string     orders results by 'name', 'created' or 'lastModified' (default "lastModified")
      --sortOrder string   'asc' or 'desc'
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
