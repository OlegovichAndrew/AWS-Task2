# AWS-Task2

There are two applications.  
The logic is next:

* HTTP server if trigered by its endpoint send request to gRPC server, then receive the file  
and increase its value if it's a text file, finally push it back to S3 bucket (like another version of existing one).
* gRPC server downloads the file from AWS S3 bucket and send it back by gRPC.  