pipeline {
    agent { 
        label "docker"
     }
    environment {
      DOCKER = credentials('docker_hub')
    }
    stages {
        stage('Build') {
            parallel {
                stage('Express Image') {
                    steps {
                        sh 'docker build -f Dockerfile \
                        -t ekas-portal-api-dev:latest .'
                    }
                }                    
            }
            post {
                failure {
                    echo 'This build has failed. See logs for details.'
                }
            }
        }
        stage('Test') {
            steps {
                echo 'This is the Testing Stage'
            }
        }
        stage('Deploy') {
            steps {
                echo 'This is the Deploy Stage'
            }
        }
    }
}