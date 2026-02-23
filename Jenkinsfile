pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage('Build') {
            steps {
                setGitHubPullRequestStatus(
                    context: 'ci/jenkins/build',
                    message: 'Building the Application',
                    state: 'PENDING'
                )

                sh 'echo Build Stage'
            }

            post {
                success {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/build',
                        message: 'Build finished successful',
                        state: 'SUCCESS'
                    )
                }
                failure {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/build',
                        message: 'Build failed',
                        state: 'FAILURE'
                    )
                }
            }
        }

        stage('Security Scan') {
            steps {
                setGitHubPullRequestStatus(
                    context: 'ci/jenkins/security',
                    message: 'Scanning for vulnerabilities...',
                    state: 'PENDING'
                )

                sh 'echo Security Scan Stage'
            }
            post {
                success {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/security',
                        message: 'No critical vulnerabilities found.',
                        state: 'SUCCESS'
                    )
                }
                failure {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/security',
                        message: 'Security vulnerabilities detected!',
                        state: 'FAILURE'
                    )
                }
            }
        }

        // --- 3. Stage: Unit Test ---
        stage('Unit Test') {
            steps {
                setGitHubPullRequestStatus(
                    context: 'ci/jenkins/test',
                    message: 'Running unit tests...',
                    state: 'PENDING'
                )

                sh 'echo Unit Test Stage'
            }
            post {
                success {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/test',
                        message: 'All tests passed!',
                        state: 'SUCCESS'
                    )
                }
                failure {
                    setGitHubPullRequestStatus(
                        context: 'ci/jenkins/test',
                        message: 'Some tests failed.',
                        state: 'FAILURE'
                    )
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
