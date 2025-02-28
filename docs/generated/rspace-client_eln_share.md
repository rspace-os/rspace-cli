## rspace-client eln share

Shares one or more documents or notebooks with one or more users or groups

### Synopsis

Share documents or notebooks with groups and individual users
	Documents/notebooks can be specified by plain IDs (e.g. 12345) or by Global Ids (e.g. SD12345)
	

```
rspace-client eln share [flags]
```

### Examples

```

// share a document and notebook with edit permission with a group into a designated folder
rspace eln share SD12345 NB23456 --groups 122345--permission edit --folder 7689

// Share a document with users in my group with read permission
rspace eln share 12345 --users 122345,45678


	
```

### Options

```
      --folder int          Target folder id (group sharing only)
      --groups string       Comma separated list of group Ids to share with
  -h, --help                help for share
      --permission string   Permission - 'read' or 'edit' (default "read")
      --users string        Comma separated list of user Ids to share with
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
