pipeline {
   agent any
   environment {
       dockerHub = "docker.io"
       docker_cred = 'dockerhub'
       http_server_path = "./AWS_HTTP-server"
       grpc_server_path = "./AWS_gRPC-server"

   }
   stages {

		stage("Build Docker Image for HTTP-server") {
			steps{
				sh "docker build -t andrewolegovich/http-server ${http_server_path}"
			}
		}

        stage("Build Docker Image for gRPC-server") {
            steps{
                sh "docker build -t andrewolegovich/grpc-server ${grpc_server_path}"
            }
        }

		stage('Upload Images to DockerHub') {

			steps {
	     	    withCredentials([
     		 	[$class: 'UsernamePasswordMultiBinding', credentialsId: docker_cred, usernameVariable: 'dockeruser', passwordVariable: 'dockerpass'],
  				]){
					sh "docker login -u ${dockeruser} -p ${dockerpass} ${dockerHub}"
  				}
	    	  	sh "docker push andrewolegovich/grpc-server"
	    	  	sh "docker push andrewolegovich/http-server"
	    	 }

	  	}
   }
}