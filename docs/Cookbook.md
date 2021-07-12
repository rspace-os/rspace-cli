Here are some ideas on how to use the CLI for various use-cases.

## 1. Uploading a folder of files to the Gallery

### Scenario
You have many files to upload and doing it through the web interface is tedious.

### Solution
    rspace eln upload myFolder/ --recursive --add-summary

This will upload all files in directory (and its subdirectories) and upload them, generating a summary document with links to all the uploaded files.

If you want more control over what you upload, you can combine with standard Unix tools, e.g.

    find . -name "*.jpg" | xargs -I file rspace eln upload file --add-summary

All the files uploaded will be placed in the 'ApiInbox' folder of their respective Galleries.

## 2. Downloading a folder of files  from the Gallery

### Scenario

You have a folder of files in RSpace Gallery and you want to download them all. Using the web interface to download each file one at a time is tiresome.

### Solution

Use the `listTree` command to navigate through   file tree. Or, you can get the folder ID from the 'info' popup in the web application.

    rspace eln listTree 

Once you have a folder ID (in this example, 9), get the folder contents in 'quiet' mode. This just lists the IDs, which can then be piped into a second command to download the files

    rspace eln listTree --folder 9 -f quiet | xargs ./rs eln download

## 3. Sharing many items at once

### Scenario

You have a folder full of RSpace documents. You'd like to share them all with your group, but don't want them organised into a notebook, and sharing them one a time is tedious and error-prone

### Solution

To share 1 or more items we can use the `share` command. To share a single item:

    rspace eln share SD2345 --groups 123 --folder 567

will share document to the group with id '123' to folder with id '567'.


To share many items, we start with using `listTree` to identify the folder with the documents we want to share.

Now, we call `groups` to find the id of the group we want to share with, and the id of the group shared folder

    rspace eln groups

```
    Id      	Name                     	Type      	SharedFolderId  
    7       	g1                       	LAB_GROUP 	326
```

By listing the items to share in quiet mode, we can pipe the ids to the share command:

    rspace eln listTree --filter document,notebook --folder 2345 -f quiet | xargs -I item  share item --groups 7 --folder 326 --permission edit

Note that we include the `--filter` option - this is so that we don't attempt to share any subfolders, which is not possible in RSpace. 

If we wanted to share with 'read' permission only we can omit the `--permission` argument, or explicitly state `--permission read`.

## 4. Importing many Word / OpenOffice Documents

### Scenario

You have  many Word documents, and you'd like to import them into RSpace as native RSpace documents. 

### Solution

This works in a similar way to `upload`. You can supply one or more folders and files, and 
the command will filter for Word documents based on extension (.odt, .doc, .docx), then import them into RSpace.

    rspace eln importWord myfolder

If you're not quite sure what you'll end up importing, you can use the `--dry-run` flag which will show you what would be uploaded, but doesn't change anything on RSpace.

    rspace eln importWord myfolder --dry-run

## 5. Inspecting XML exports 

### Scenario 

You have been making regular XML exports of your RSpace documents and have accumulated many .zip files over time, some of them quite big. You'd like to know what's inside withoout having to 
unzip or import back into RSpace.

### Solution

The `archive` command is different from other commands in this software as it doesn't make calls to an RSpace server - it works with files on your device. 

    rspace archive myArchive1.zip myArchive2.zip --summary

will parse the archive (without extracting it), print out the manifest file and summarise the content of the archive, including the date range of documents created in the archive and also
a list of authors.

Example output:

```
file      	Total Docs	minDate               	maxDate               	Authors                                           
rs2.zip   	2         	2020-05-10T18:39:05Z  	2020-05-16T09:21:24Z  	user5e                                            
rs3.zip   	3         	2020-05-02T11:32:09Z  	2020-05-17T20:01:58Z  	user5e;bobsmith
```

If you want to find out more information about a single archive, you can use the extended summary flag `--xsummary`; this will list the names of documents in the archive

```
    rspace archive myArchive1.zip myArchive2.zip --xsummary
```

## 6. Creating partially filled content automatically for Structured (multi-field) documents

### Scenario

You are using RSpace Forms to capture structured data from PCR experiments. Some of the data records the PCR setup, and the source of this data can be obtained in a spreadsheet format. You'd like to pre-populate a set of documents with these setups

### Solution

The `addDocument` command can create multiple documents from structured data.

1. First of all create an RSpace form to record the PCR experiment. This might be number of cycles (number), denaturing time (time in seconds), annealing time (time in seconds), extension time (time in seconds), 5prime oligo (string) , 3prime oligo (string), Notes (text) and Results (text). 
Note the ID of the form -  let's say it's FM12345

2. Create/export a set of PCR experiment setups in CSV format, e.g.

```
    cycles,denature,annealing,extension,5prime,3prime,Notes,Results
    35,300,45,120,atgctagcgctagc,atgcacgggcacac,,
    30,270,40,150,atcgagctagtc,catcgctacgtcg,,
```
Each row maps to an RSpace document; each column maps to a field in the form.


Note the last 2 columns are left blank - this is for manual description of the results which will be added in the web application later. Save this CSV data in a file 'myPcrSetup.csv'.

3. Create a new notebook to hold these experiments and note the ID in the output - let's assume its NB678

```
rspace eln addNotebook --name myPcrExperiments 
```

4. Now, in a single command, you can create the documents automatically with your experimental setup pre-populated:

```
rspace eln addDocument --formId FM12345 --name myPcrExperiment --input myPcrSetup.csv --folder NB678
```
 
## 7. Exporting to XML  and HTML

### Scenario

You'd like to make regular exports of a folder, or a set of search results.
You could do this in the web application, but you'd have to remember to do the export,
wait for it to complete and download somewhere. 

### Solution

The `export` command is what you need here. You can export a selection, all your work or all
your group's work. Here we'll focus on exporting a selection.

If you have a folder or a notebook, with lots of content inside, then you can just export the folder - all child content will be included in the export. You can get the ID from Workspace listing or from the 'GetInfo' page. Let's suppose you have a folder FL12345. You can use the 'global id' (FL12345) or just the numeric ID (12345), whatever is more convenient.

The `--wait` option blocks till the export is complete.

```
rspace eln export 12345 --scope selection --format html --wait
```

If you have a more complex selection - say a bunch of documents that share a common tag, you 
can combine export with search. E.g.

```
rspace eln listDocuments --tag grantNumber1234 --maxResults 100  -f quiet | \
   xargs rspace eln export --scope selection --format html --wait
```

Here we do a search with the 'quiet' flag set which just outputs the IDs of hits.
Using `xargs` we can pipe these IDs into our export command. 

Once export is completed, you get a job ID which you can use to download the results.

The CLI doesn't automatically download the output zip file as it might be really huge. The output of `export` tells you the size:

```
Id      	Status    	Percent Complete  	Download size 
91      	COMPLETED 	100.00            	949 MB 
```

You can download using the `job` command. Don't forget the `--download` flag, else it will
just show the job status.

```
rspace eln job 91 --download 
```

If you really want to export and download in one go, you can use `xargs` again.
Here we export, wait, and download a whole  user's work  in a single line.

```
rspace eln  export 123 --scope user  --wait -f quiet | \
   xargs -I jobId  rspace eln  job  jobId  --download
```

This latter command could be used as an input to  `cron`. What you do from here is up to you - send to a long-term archive or repository, send to collaborators etc.

## 8.  Exporting many users's work

### Scenario

You have many users' work to export, perhaps to transfer to an archive or to transfer from another RSpace server

### Solution

First of all get the usernames and ids of user:

    ./rscli eln listUsers --maxResults 100 -f csv | grep -v '-' | cut -d ',' -f1,2 > users.csv

This gets users, filters out the relevant fields and saves them in a CSV file. (Note there is a limit of 100 users currently).

Now iterate over this list, export and download:

```
   logFile=$(date "+%d-%m-%y-%T").log
   outfileDir=/media/rspace/exports
   echo "Logging progress to $logFile"
   while read person; do
      id=$(echo $person |  awk -F "," '{print $1}')
      user=$(echo $person |  awk -F "," '{print $2}')
      echo "exporting $user"
      echo "exporting $user" >> $logFile
      jobid=$(./rscli eln export $id --scope user --format xml --wait | tail -n +2 | cut -f1)
      echo "downloading export from jobID $jobid"
      ./rscli eln job $jobid --download --outfile ${outfileDir}/${user}.zip
    done < users.csv
echo "Completed"

