pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage('Clean work space') {
            steps {
                cleanWs()
            }
        }

        stage('Build') {
            steps {
                sh 'echo Starting Build Stage...'

                sh 'echo Build Stage'
            }

            post {
                success {
                    sh 'echo Build finished successfully!'
                }
                failure {
                    sh 'echo Build failed!'
                }
            }
        }

        stage('Security Scan') {
            steps {
                sh 'echo Scanning for vulnerabilities...'

                sh 'echo Security Scan Stage'
            }
            post {
                success {
                    sh 'echo No critical vulnerabilities found.'
                }
                failure {
                    sh 'echo Security vulnerabilities detected!'
                }
            }
        }

        // --- 3. Stage: Unit Test ---
        stage('Unit Test') {
            steps {
                sh 'echo Running unit tests...'

                sh 'echo Unit Test Stage'
            }
            post {
                success {
                    sh 'echo All tests passed!'
                }
                failure {
                    sh 'echo Some tests failed.'
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
