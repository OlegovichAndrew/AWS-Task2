# AWS-Task2

There are two applications.  
The logic is next:

* gRPC server downloads the file from AWS S3 bucket and send it back by gRPC.  
* HTTP server if trigered by endpoint send request to gRPC server, then receive the file  
and increase its value if it's a text file, then push it back to S3 bucket (like another version of existing one).