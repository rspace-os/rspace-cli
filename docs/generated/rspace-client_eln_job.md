## rspace-client eln job

Query progress of a Job

### Synopsis

 Query a job status, or download result. You get a jobId after submitting an export
	 request using the 'export' command.
	

```
rspace-client eln job [flags]
```

### Examples

```

// get progress of job in tabular format
rspace eln job  22

// get raw JSON
rspace eln job 22 -f json

// download (if complete) to a file of your choice
rspace eln job  22 --download --outfile /home/me/myexport.zip

// download (if complete) to current directory
rspace eln job  22 --download
	
```

### Options

```
      --download         Download result, if complete
  -h, --help             help for job
      --outfile string   file path to download to...
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
