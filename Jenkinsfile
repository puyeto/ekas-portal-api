pipeline {
    environment {
      DOCKER = credentials('docker_hub')
    }
    agent any
        stages {
            stage('Build') {
                parallel {
                    stage('Express Image') {
                        steps {
                            sh 'docker build -f express-image/Dockerfile \
                            -t ekas-portal-apiapp-dev:trunk .'
                        }
                    }
                    stage('Test-Unit Image') {
                        steps {
                            sh 'docker build -f test-image/Dockerfile \
                            -t test-image:latest .'
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