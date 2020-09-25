#
pipeline {
    agent any

    triggers { cron("0 * * * *") }


    parameters {
        string(name: 'DB_BRANCH', defaultValue: 'master', description: 'What branch of the pipeline Dashboard should be deployed?')
        string(name: 'BRANCH', defaultValue: 'master', description: 'What branch of the infra-infra config should be deployed?')
    }


    environment {
        GOPATH = "${WORKSPACE}/go"
        GOCACHE = "${WORKSPACE}/.go-build"
        NPM_CONFIG_CACHE = "${WORKSPACE}/.npm"
        CI = true
        KUBECONFIG = credentials('ci_dashboard_deploy_user')
        DOCKER_CONFIG="${WORKSPACE}/.docker"
    }

    stages {
        stage('Test') {
            steps {
              node('worker:artifactory.delivery.puppetlabs.net/dev-services/node-go-java') {

            }
        }
        stage('Build') {
            steps {
              node('worker:artifactory.delivery.puppetlabs.net/dev-services/node-go-java') {
                
            }
        }
        stage('Update') {
            steps {
            }
        }
    }
}
