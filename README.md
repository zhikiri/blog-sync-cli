# Blog synchronization tool

Simple CLI util for synchronize [hugo](https://gohugo.io) static content with [AWS S3](https://aws.amazon.com/s3/) bucket.

## Description

CLI tool sync new/changed/removed files. Currently using checksum tests.

Some description notes:

- Tool will ask required information (AWS credentials, path to files, ignore extentions list) and store in hidden configuration file (ex `~/.bsync/setup.json`)

- Add ignore extentions list for avoid synchronization for some files
