## RSpace-CLI

RSpace-CLI is a command-line application to interact with [RSpace ELN](https://www.researchspace.com) in a more convenient
and compact way than using the API directly. It supports scripting and programmatic access
to RSpace for automation, bulk operations and routine functions and can be used independently or in conjunction with RSpace web application.

It is designed a supplement to the web interface for tasks such as:

* Bulk upload or download of files and folders
* Bulk import of MSWord documents into native RSpace documents
* Querying audit-trail for activity
* Getting reports in JSON, tabular or CSV format
* Integrating cleanly into your data-management workflows.
* Admin functions such as ad-hoc account creation.
* Creating notebook entries pre-populated with content.
* Inspecting contents of zipped exports.

A [Cookbook](docs/Cookbook.md) shows some uses of the CLI.

It is written in the Go programming language.

### Downloading

Binaries are available as releases on this Github project site. 

Once downloaded, make the file executable and optionally create a convenient alias, e.g. 

    ln -s my-download rspace
    chmod 755 rspace
    ./rspace eln --help

to show commands and their arguments. You can also view an [HTML version](docs/generated/rspace-client_eln.md) of the documentation.

### Configuring

Next, you must supply a configuration file with your RSpace API credentials:

Create a file called '.rspace' in your home folder and add two lines with the URL of your RSpace and
your API key, like this:

    RSPACE_API_KEY=get_this_from_your_RSpace_profile_page
    RSPACE_URL=https://myrspace.com/api/v1

If you prefer, instead of the default '.rspace' file,  you can add this information to any file, save it with a '.env' suffix and supply its filepath with the --config flag to each command, e.g.

    rspace eln listTree --config /path/to/myConfig.env

Using --config option is useful if you have more than one account (e.g. an admin account and a personal account)
