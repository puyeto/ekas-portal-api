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
                            sh 'docker build . -f Dockerfile \
                            -t omollo/ekas-portal-api-prod:latest .'
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
            stage('DEPLOY') {
                when {
                    branch 'master'  //only run these steps on the master branch
                }
                steps {
                    // sh 'docker tag ekas-portal-api-dev:latest omollo/ekas-portal-api-prod:latest'
                    sh 'docker login -u "omollo" -p "safcom2012" docker.io'
                    sh 'docker push omollo/ekas-portal-api-prod:latest'
                }
            }

            stage('PUBLISH') {
                when {
                    branch 'master'  //only run these steps on the master branch
                }
                steps {
                    // sh 'docker swarm leave -f'
                    // sh 'docker run -d -p 8081:8081 --rm --name ekas-portal ekas-portal-api-dev'
                    // sh 'docker swarm init --advertise-addr 159.89.134.228'
                    sh 'docker stack deploy -c docker-compose.yml ekas-portal-api-prod'
                }

            }

            // stage('REPORTS') {
            //     steps {
            //         junit 'reports.xml'
            //         archiveArtifacts(artifacts: 'reports.xml', allowEmptyArchive: true)
            //         // archiveArtifacts(artifacts: 'ekas-portal-api-prod-golden.tar.gz', allowEmptyArchive: true)
            //     }
            // }

            stage('CLEAN-UP') {
                steps {
                    // sh 'docker stop ekas-portal-api-dev'
                    sh 'docker system prune -f'
                    deleteDir()
                }
            }
        }
    }