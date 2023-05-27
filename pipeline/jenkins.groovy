pipeline {
    agent any
    parameters {
        choice(name: 'OS', choices: ['linux', 'darwin', 'windows', 'all'], description: 'Pick OS')
        choice(name: 'ARCH', choices: ['amd64', 'arm64'], description: 'Pick ARCH')
    }
    environment {
        REPO = "https://github.com/xminyax/kbot.git"
        BRANCH = 'main'
    }
    stages {
        stage('Example') {
            steps {
                echo "Build for platform ${params.OS}"
                echo "Build for arch: ${params.ARCH}"
            }
        }
        stage("clone") {
            steps {
                echo 'CLONE REPOSITORY'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }
        stage("test") {
            steps {
                sh "make --version"
                echo 'TEST EXECUTION STARTED '
                sh 'make test'
            }
        }
        stage("build") {
            steps {
                echo 'BUILD EXECUTION STARTED '
                sh 'make build'
            }
        }
        stage("image") {
            steps {
                script {
                    echo 'TEST EXECUTION STARTED '
                    sh 'make image'
                }
            }
        }
        stage("push") {
            steps {
                script {
                    docker.withRegistry('', 'dockerhub') {
                        sh 'make push'
                    }
                }
            }
        }
    }
}