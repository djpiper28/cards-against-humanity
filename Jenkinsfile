pipeline {
    agent any

    stages {
        stage('Build and package backend') {
          steps {
            sh 'docker build -t cahbackend .'
            sh 'docker tag cahbackend localhost:5000/cahbackend'
            sh 'docker push localhost:5000/cahbackend'
          }
        }
        
        stage('Package') {
          steps {
            sh '/var/lib/jenkins/go/bin/helmVersioner charts/cahbackend/Chart.yaml'
            sh 'helm install cahbackend charts/cahbackend || true'
            sh 'helm upgrade cahbackend charts/cahbackend'
          }
        }

        stage('Build Frontend') {
            steps {
              sh 'cd cahfrontend && pnpm i'
              sh 'cd cahfrontend && pnpm run build'
          }
        }

        stage('Deploy') {
          steps {
            sh 'rm -r /home/static/cah/* || true'
            sh 'cp -r cahfrontend/dist/* /home/static/cah/'
            sh 'sudo /home/scripts/restart_nginx'
          }
        }
    }
}
