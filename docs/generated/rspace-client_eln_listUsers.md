## rspace-client eln listUsers

Lists users - requires sysadmin role!

### Synopsis

List users, sorted or paginated.

Please note that currently users are ordered by account creation date only (default is most recent first),
but there is no limit on the page size, so you can retrieve all users in one query, then
filter or sort using standard Unix utilities.

	

```
rspace-client eln listUsers [flags]
```

### Examples

```

// find out how many users you have:
rspace eln listUsers --all -f json | jq '.TotalHits'

// or just use 'grep' if you don't have 'jq' installed:
rspace eln listUsers --all -f json | grep TotalHits

// now you can get all users in one go and sort by any column e.g. by username:
rspace eln listUsers --all | sort -k2
	
```

### Options

```
      --all                Gets all users
  -h, --help               help for listUsers
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
