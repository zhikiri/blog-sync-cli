# Blog sync CLI

Command line tool that provide synchronization mechanism for my blog.
Basic concept of the library is files checksum comparison, it will upload only modified files or delete them if it's required.

## Configuration

Before start to use it, need to create a configuration file from example (`config.json.example`).
Also, two environment variables are required (`AWS_ACCESS_KEY` and `AWS_ACCESS_SECRET`).

## Conclusion

In current version of the library there is a possibility to have different synchronization places, not just only s3 bucket.
Maybe in future it can be interesting thing to do.
