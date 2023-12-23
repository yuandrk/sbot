pipeline {
    agent any
    environment {
        REPO = 'https://github.com/yuandrk/sbot'
        BRANCH = 'main'
        REGISTRY = 'ghcr.io/yuandrk'
        APP = 'sbot'
        GHCR = 'https://ghcr.io'

    }
    parameters {
        choice(name: 'OS', choices: ['linux', 'darwin', 'windows', 'all'], description: 'Choose OS')
        choice(name: 'ARCH', choices: ['amd64', 'arm64'], description: 'Choose architectory')

    }
    stages {
        stage("clone") {
            steps {
                echo 'ClONE REPOSITORY'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }

        stage("test") {
            steps {
                echo 'TEST STARTED'
                sh 'make test'
            }
        }

        stage("build") {
            steps {
                echo 'BUILD STARTED'
                sh 'make build TARGETOS=${OS} TARGETARCH=${ARCH}'
            }
        }

        stage("image") {
            steps {
                echo 'IMAGE STARTED'
                sh 'make image REGISTRY=${REGISTRY} APP=${APP}'
            }
        }
        stage("login to GHCR") {
            steps {
                echo 'LOGIN TO GHCR'
                withCredentials([usernamePassword(credentialsId: 'github-token', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                    sh 'docker login -u $USERNAME -p $PASSWORD ${GHCR}'
                }

            }
        }

        stage("push") {
            steps {
                echo 'PUSH TO GHCR'
                sh 'make push REGISTRY=${REGISTRY} TARGETOS=${OS} TARGETARCH=${ARCH}'
            }
        }
    }
}