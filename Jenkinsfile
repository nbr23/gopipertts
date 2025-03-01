pipeline {
    agent any

    options {
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout'){
            steps {
                checkout scm
            }
        }
        stage('Prep buildx') {
            steps {
                script {
                    env.BUILDX_BUILDER = getBuildxBuilder();
                }
            }
        }
        stage('Build docker Image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKERHUB_CREDENTIALS_USR', passwordVariable: 'DOCKERHUB_CREDENTIALS_PSW')]) {
                    sh 'docker login -u $DOCKERHUB_CREDENTIALS_USR -p "$DOCKERHUB_CREDENTIALS_PSW"'
                }
                sh """
                    docker buildx build \
                        --pull \
                        --builder \$BUILDX_BUILDER  \
                        --platform linux/amd64,linux/arm64,linux/arm/v7 \
                        -t nbr23/gopipertts:latest \
                        ${ "$GIT_BRANCH" == "master" ? "--push" : ""} .
                    """
            }
        }
        stage('Sync github repo') {
            when { branch 'master' }
            steps {
                syncRemoteBranch('git@github.com:nbr23/gopipertts.git', 'master')
            }
        }
    }
    post {
        always {
            sh 'docker buildx stop $BUILDX_BUILDER || true'
            sh 'docker buildx rm $BUILDX_BUILDER || true'
        }
    }
}
