pipeline {
    agent any

    tools{
        nodejs 'Node25'
    }
    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 1, unit: 'HOURS')
        disableConcurrentBuilds()
    }

    environment {
        SCANNER_HOME=tool 'sonar-scanner'
        DOCKER_IMAGE="0067chetsada/medee-app"
        IMAGE_TAG="${env.BUILD_NUMBER}"
        SONAR_TOKEN=credentials('Sonar-token')
        GITOPS_BRANCH = "main"
    }

    stages {
        stage('Clean work space') {
            steps {
                cleanWs()

                checkout scm
            }
        }

        stage('Frontend test'){
            steps {
                script {
                    docker.withRegistry('https://dhi.io', 'Docker') {
                        sh 'DOCKER_BUILDKIT=1 docker build --target run-test-stage -f frontend/Dockerfile frontend/'
                    }
                }
            }
        }

        stage('Backend test'){
            steps {
                script {
                    docker.withRegistry('https://dhi.io', 'Docker') {
                        sh """
                            cd backend 
                            go mod download
                            go test -v ./...
                            """
                    }
                }
            }
        }

        stage('Dependency Scan'){
            steps {
                echo'Frontend Depency Install'
                sh """
                    cd frontend && \\
                    npm install --global corepack@latest && \\
                    corepack enable && \\
                    corepack prepare pnpm@latest-10 --activate && \\
                    pnpm install
                    """

                echo 'Dependency Scan'
                dependencyCheck additionalArguments: '--scan ./ --disableYarnAudit --disableNodeAudit', odcInstallation: 'DP'
                dependencyCheckPublisher pattern: '**/dependency-check-report.xml'
            }
        }

        stage('Security Scan'){
            steps {
                sh 'trivy scan2html fs --cache-dir /var/lib/jenkins/.cache . --scan2html-flags --output trivyfs-report.html'

                archiveArtifacts artifacts: 'trivyfs-report.html', allowEmptyArchive: true
            }
        }

        stage("Security & Code Analysis") {
            steps {
                withSonarQubeEnv('sonar-scanner') {
                    sh """
                    $SCANNER_HOME/bin/sonar-scanner \\
                    -Dsonar.projectName=Medee \\
                    -Dsonar.projectKey=Medee \\
                    -Dsonar.token=${SONAR_TOKEN} \\
                    """
                }
            }
        }

    stage ("Quality Gate") {

        // steps {
        //     waitForQualityGate abortPipeline: true
        // }

        steps {
            // set true when don't want to continue to next test if not passing quality
            waitForQualityGate abortPipeline: false, credentialsId: 'Sonar-token'
        }
        // timeout(time: 1, unit: 'HOURS'){
        //     def qg = waitForQualityGate()
        //     if(qg.status != 'OK') {
        //         error "Pipeline aborted dueto quality gate failure: ${qg.status}"
        //     }
        // }
    }

        stage('Build & Image Scan Backend'){
                steps{
                    script {
                        docker.withRegistry('https://dhi.io', 'Docker') {
                            sh """
                                    cd backend && \\
                                    docker build -t ${DOCKER_IMAGE}-backend:${env.BUILD_NUMBER} . &&\\
                                    trivy image ${DOCKER_IMAGE}-backend:${env.BUILD_NUMBER}
                                """
                        }
                    }
                }
        }
        
        stage('Build & Image Scan Frontend'){
            steps{
                script {
                    docker.withRegistry('https://dhi.io', 'Docker') {
                        sh """
                                cd frontend && \\
                                docker build -t ${DOCKER_IMAGE}-frontend:${env.BUILD_NUMBER} . && \\
                                trivy image ${DOCKER_IMAGE}-frontend:${env.BUILD_NUMBER}
                            """
                    }
                }
            }
        }

        stage('Push to Registry') {
            steps {
                script {
                    docker.withDockerRegistry(credentialsId: 'Docker', toolName: 'docker') {
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

        stage('Update GitOps Repository') {
            steps {
                script {
                    withCredentials([string(credentialsId: 'medee-gitops-github-app', variable: 'GITHUB_TOKEN')]) {
                        sh """
                            # Clone GitOps repo with token
                            rm -rf medee-gitops
                            git clone https://x-access-token:${GITHUB_TOKEN}@github.com/Nut06/medee-gitops.git
                            cd medee-gitops
                            
                            # Configure git
                            git config user.email "jenkins@medee.local"
                            git config user.name "Jenkins CI"
                            
                            # Update backend image tag
                            sed -i 's|image: ${DOCKER_IMAGE}-backend:.*|image: ${DOCKER_IMAGE}-backend:${IMAGE_TAG}|g' apps-config/medee-dev/backend/deployment.yaml
                            
                            # Update frontend image tag
                            sed -i 's|image: ${DOCKER_IMAGE}-frontend:.*|image: ${DOCKER_IMAGE}-frontend:${IMAGE_TAG}|g' apps-config/medee-dev/frontend/deployment.yaml
                            
                            # Commit and push
                            git add apps-config/medee-dev/backend/deployment.yaml
                            git add apps-config/medee-dev/frontend/deployment.yaml
                            git commit -m "chore: update images to build ${IMAGE_TAG}" || echo "No changes to commit"
                            git push origin ${GITOPS_BRANCH}
                            
                            # Cleanup
                            cd ..
                            rm -rf medee-gitops
                        """
                    }
                }
            }
}


        // --- 4. GitOps (ArgoCD Trigger) ---
        stage('Trigger GitOps (ArgoCD)') {
            steps {
                echo "Deploying Medee Version: ${IMAGE_TAG} to EKS via ArgoCD..."
                echo "ArgoCD will auto-sync the changes from GitOps repo"
            }
        }

    }

    post {
        
        cleanup {
            /* clean up our workspace */
            deleteDir()
            /* clean up tmp directory */
            dir("${workspace}@tmp") {
                deleteDir()
            }
            /* clean up script directory */
            dir("${workspace}@script") {
                deleteDir()
            }
        }
        // sending email & remove build docker image
        always {
            sh 'docker image prune -f || true'
            sh 'docker system df'

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
                attachmentsPattern: 'trivyfs-report.html,trivyimage.txt'
            )
            }
        }
    }
}
