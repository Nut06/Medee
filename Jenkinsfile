pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 1, unit: 'HOURS')
        disableConcurrentBuilds()
    }

    environment {
        SCANNER_HOME=tool 'sonar-scanner'
        DOCKER_IMAGE="chetsada/medee-app"
        IMAGE_TAG="${env.BUILD_NUMBER}"
        SONAR_TOKEN=credentials('Sonar-token')
    }

    stages {
        stage('Clean work space') {
            steps {
                cleanWs()

                checkout scm
            }
        }

        stage('Automate test'){
            parallel{
                stage('Frontend (Vitest)') {
                    steps {
                        script{
                            sh 'docker build --target run-test-stage -f frontend/Dockerfile frontend/'
                         }
                    }
                }

                stage('Backend (Go)') {
                    steps {
                        script{
                            sh 'docker build --target run-test-stage -f backend/Dockerfile backend/'
                        }
                    }
                }
            }
        }

        stage('Security Scan'){
            steps {
                sh 'trivy fs . > trivyfs-report.txt'
            }
        }

    stage("Security Analysis"){
            parallel {
                stage('SonarQube Analysis') {
                    steps {
                        withSonarQubeEnv('sonar-server') {
                            sh """
                            $SCANNER_HOME/bin/sonar-scanner \\
                            -Dsonar.projectName=Medee \\
                            -Dsonar.projectKey=Medee \\
                            -Dsonar.branch.name=${env.BRANCH_NAME} \\
                            -Dsonar.token=${SONAR_TOKEN} \\
                            """
                        }
                    }
                }

                stage('OWASP & Trivy FS Scan') {
                    steps {
                        // สแกน Library (SCA)
                        dependencyCheck additionalArguments: '--scan ./ --disableYarnAudit --disableNodeAudit', odcInstallation: 'DP'
                        dependencyCheckPublisher pattern: '**/dependency-check-report.xml'
                        
                        // สแกน File System
                        sh "trivy fs . > trivyfs-report.txt"
                    }
                }
            }
    }

    stage ("Quality Gate") {
        steps {
            waitForQualityGate abortPipeline: true, credentialsId: 'Sonar-token'
        }
    }

        stage('Docker Build & Image Scan'){
            parallel {
                stage('Build & Scan Backend'){
                    steps{
                        script{
                            sh 'cd backend'
                            sh 'docker build -t ${DOCKER_IMAGE}-backend:${env.BUILD_NUMBER} .'

                            sh 'trivy image ${DOCKER_IMAGE}-backend:${env.BUILD_NUMBER}'
                        }
                    }
                }

                stage('Build & Scan Frontend'){
                    steps{
                        script{
                            sh 'cd frontend'
                            sh 'docker build -t ${DOCKER_IMAGE}-frontend:${env.BUILD_NUMBER} .'

                            sh 'trivy image ${DOCKER_IMAGE}-frontend:${env.BUILD_NUMBER}'
                        }
                    }
                }
            }
        }

        stage('Push to Registry') {
            steps {
                script {
                    docker.withDockerRegistry('https://index.docker.io/v1/', 'Docker') {
                        // Push ทั้งเวอร์ชันระบุเลข Build และเวอร์ชัน latest
                        sh "docker push ${DOCKER_IMAGE}-frontend:${IMAGE_TAG}"
                        sh "docker push ${DOCKER_IMAGE}-backend:${IMAGE_TAG}"
                        sh "docker tag ${DOCKER_IMAGE}-frontend:${IMAGE_TAG} ${DOCKER_IMAGE}-frontend:latest"
                        sh "docker tag ${DOCKER_IMAGE}-backend:${IMAGE_TAG} ${DOCKER_IMAGE}-backend:latest"
                        sh "docker push ${DOCKER_IMAGE}-frontend:latest"
                        sh "docker push ${DOCKER_IMAGE}-backend:latest"
                        app.push("latest")
                    }
                }
            }
        }

        // --- 4. GitOps (ArgoCD Trigger) ---
        stage('Trigger GitOps (ArgoCD)') {
            steps {
                echo "Deploying Medee Version: ${IMAGE_TAG} to EKS via ArgoCD..."
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
