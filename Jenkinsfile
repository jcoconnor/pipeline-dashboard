#
pipeline {
    agent any

    triggers { cron("0 * * * *") }

    parameters {
        string(name: 'version', defaultValue: '', description: 'This will be used to set the version number and the docker image tag, Example 0.1.2')
    }

    environment {
        GOPATH = "${WORKSPACE}/go"
        GOCACHE = "${WORKSPACE}/.go-build"
        NPM_CONFIG_CACHE = "${WORKSPACE}/.npm"
        CI = true
        DOCKER_CONFIG="${WORKSPACE}/.docker"
    }

    stages {
        stage('Build') {
            steps {
              node('worker:artifactory.delivery.puppetlabs.net/dev-services/node-go-java') {
                  sh './build.sh'
            }
        }
    }
}
