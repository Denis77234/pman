

## The package manager provides the following commands:

    pman create ./packet.json: Packages the files specified in the package file into an archive and send to the server via SSH.
    pman update ./packages.json: Downloads archive files via SSH.


    ## Package File Format
The package file should have `.json` format. It should include paths to select files using glob patterns.

## Example Package File:
**packet.json**

```json
{
  "name": "packet-1",
  "ver": "1.10",
  "targets": [
    "./archivethis1/*.txt",
    {"path": "./archivethis2/", "exclude": "*.tmp"}
  ],
  "packets": [
    {"name": "packet-3", "ver": "<=2.0"}
  ]
}

