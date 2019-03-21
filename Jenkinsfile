pipeline {
    environment {
      DOCKER = credentials('docker-hub')
    }
    agent any
        // Building Images
        stages {
            stage('Build') {
                parallel {
                    stage('Express Image') {
                        steps {
                            sh 'docker build -f express-image/Dockerfile \
                            -t ekas-portal-api-dev:trunk .'
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
            // Performing Software Tests
            stage('Test') {
                steps {
                    echo 'This is the Testing Stage'
                }
            }
            // Deploy Software
            stage('Deploy') {
                when {
                    branch 'master'  //only run these steps on the master branch
                }
                steps {
                    retry(3) {
                        timeout(time:10, unit: 'MINUTES') {
                            sh 'docker tag ekas-portal-api-dev:trunk <DockerHub Username>/ekas-portal-api-prod:latest'
                            sh 'docker push <DockerHub Username>/ekas-portal-api-prod:latest'
                            sh 'docker save <DockerHub Username>/ekas-portal-api-prod:latest | gzip > ekas-portal-api-prod-golden.tar.gz'
                        }
                    }

                }
                 post {
                    failure {
                        sh 'docker stop ekas-portal-api-dev test-image'
                        sh 'docker system prune -f'
                        deleteDir()
                    }
                }
            }

            // JUnit reports and artifacts saving
            stage('REPORTS') {
                steps {
                    junit 'reports.xml'
                    archiveArtifacts(artifacts: 'reports.xml', allowEmptyArchive: true)
                    archiveArtifacts(artifacts: 'ekas-portal-api-prod-golden.tar.gz', allowEmptyArchive: true)
                }
            }
        
            // Doing containers clean-up to avoid conflicts in future builds
            stage('CLEAN-UP') {
                steps {
                    sh 'docker stop ekas-portal-api-dev test-image'
                    sh 'docker system prune -f'
                    deleteDir()
                }
            }
            
        }
    }