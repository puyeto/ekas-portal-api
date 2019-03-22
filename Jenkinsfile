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
                when {
                    branch 'master'  //only run these steps on the master branch
                }
                steps {
                    sh 'docker tag ekas-portal-api-dev:latest omollo/ekas-portal-api-prod:latest'
                    sh 'docker login -u "omollo" -p "safcom2012" docker.io'
                    sh 'docker push omollo/ekas-portal-api-prod:latest'
                }
            }

            stage('Execute') {
                when {
                    branch 'master'  //only run these steps on the master branch
                }
                environment {
                    JENKINS_LOGIN = credentials('docker_hub')
                }
                steps {
                    echo "LOGIN=${JENKINS_LOGIN}"
                    sh 'docker run -d -p 8081:8081 -rm -name ekas-portal ekas-portal-api-dev'
                }

            }

            // stage('REPORTS') {
            //     steps {
            //         junit 'reports.xml'
            //         archiveArtifacts(artifacts: 'reports.xml', allowEmptyArchive: true)
            //         // archiveArtifacts(artifactopts: 'ekas-portal-api-prod-golden.tar.gz', allowEmptyArchive: true)
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