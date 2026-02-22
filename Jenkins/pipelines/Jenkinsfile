pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage('Build'){
            steps {
                githubNotify context: 'ci/jenkins/build', status: 'PENDING', description: 'Building the Application'

                sh 'echo Build Stage'
            }

            post {
                success {
                    githubNotify context: 'ci/jenkins/build', status: 'SUCCESS', description: 'Build finished successful'
                }
                failure {
                    githubNotify context: 'ci/jenkins/build', status: 'FAILURE', description: 'Build failed'
                }
            }
        }

        stage('Security Scan') {
            steps {
                githubNotify context: 'ci/jenkins/security', status: 'PENDING', description: 'Scanning for vulnerabilities...'
                
                sh 'echo Security Scan Stage'
            }
            post {
                success {
                    githubNotify context: 'ci/jenkins/security', status: 'SUCCESS', description: 'No critical vulnerabilities found.'
                }
                failure {
                    githubNotify context: 'ci/jenkins/security', status: 'FAILURE', description: 'Security vulnerabilities detected!'
                }
            }
        }

        // --- 3. Stage: Unit Test ---
        stage('Unit Test') {
            steps {
                githubNotify context: 'ci/jenkins/test', status: 'PENDING', description: 'Running unit tests...'
                
                sh 'echo Unit Test Stage'
            }
            post {
                success {
                    githubNotify context: 'ci/jenkins/test', status: 'SUCCESS', description: 'All tests passed!'
                }
                failure {
                    githubNotify context: 'ci/jenkins/test', status: 'FAILURE', description: 'Some tests failed.'
                }
            }
        }
    }

    post {
        // sending email
        always {
            echo 'Pipeline execution finished'
            
            script {
                def buildStatus = currentBuild.currentResult
                def buildUser = currentBuild.getBuildCauses('hudson.model.Cause$UserIdCause')[0]?.userId ?: 'Github User'

                emailext(
                subject: "Pipeline ${buildStatus}: ${env.JOB_NAME} #${env.BUILD_NUMBER}",
                body: """
                    <p>This is a Jenkins Medee CICD pipeline status.</p>
                    <p>Project: ${env.JOB_NAME}</p>
                    <p>Build Number: ${env.BUILD_NUMBER}</p>
                    <p>Build Status: ${buildStatus}</p>
                    <p>Started by: ${buildUser}</p>
                    <p>Build URL: <a href="${env.BUILD_URL}">${env.BUILD_URL}</a></p>
                """,
                // add multiple person here
                to: 'chetsadakon.chumpia@gmail.com',
                from: 'chetsadakon.chumpia@gmail.com',
                replyTo: 'chetsadakon.chumpia@gmail.com',
                mimeType: 'text/html',
                attachmentsPattern: 'trivyfs.txt,trivyimage.txt'
            )
            }
        }
    }
}
