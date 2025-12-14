// Azure Container Apps Bicep Template
// https://learn.microsoft.com/en-us/azure/container-apps/

@description('Location for all resources')
param location string = resourceGroup().location

@description('Name of the Container Apps environment')
param environmentName string = 'parse-dmarc-env'

@description('Name of the Container App')
param containerAppName string = 'parse-dmarc'

@description('IMAP Host')
@secure()
param imapHost string

@description('IMAP Username')
@secure()
param imapUsername string

@description('IMAP Password')
@secure()
param imapPassword string

resource environment 'Microsoft.App/managedEnvironments@2023-05-01' = {
  name: environmentName
  location: location
  properties: {
    zoneRedundant: false
  }
}

resource containerApp 'Microsoft.App/containerApps@2023-05-01' = {
  name: containerAppName
  location: location
  properties: {
    managedEnvironmentId: environment.id
    configuration: {
      ingress: {
        external: true
        targetPort: 8080
        transport: 'auto'
        allowInsecure: false
      }
      secrets: [
        {
          name: 'imap-host'
          value: imapHost
        }
        {
          name: 'imap-username'
          value: imapUsername
        }
        {
          name: 'imap-password'
          value: imapPassword
        }
      ]
    }
    template: {
      containers: [
        {
          name: 'parse-dmarc'
          image: 'ghcr.io/meysam81/parse-dmarc:latest'
          resources: {
            cpu: json('0.25')
            memory: '0.5Gi'
          }
          env: [
            {
              name: 'IMAP_HOST'
              secretRef: 'imap-host'
            }
            {
              name: 'IMAP_PORT'
              value: '993'
            }
            {
              name: 'IMAP_USERNAME'
              secretRef: 'imap-username'
            }
            {
              name: 'IMAP_PASSWORD'
              secretRef: 'imap-password'
            }
            {
              name: 'IMAP_MAILBOX'
              value: 'INBOX'
            }
            {
              name: 'IMAP_USE_TLS'
              value: 'true'
            }
            {
              name: 'DATABASE_PATH'
              value: '/data/db.sqlite'
            }
            {
              name: 'SERVER_PORT'
              value: '8080'
            }
            {
              name: 'SERVER_HOST'
              value: '0.0.0.0'
            }
          ]
          probes: [
            {
              type: 'Liveness'
              httpGet: {
                path: '/api/statistics'
                port: 8080
              }
              initialDelaySeconds: 10
              periodSeconds: 30
            }
            {
              type: 'Readiness'
              httpGet: {
                path: '/api/statistics'
                port: 8080
              }
              initialDelaySeconds: 5
              periodSeconds: 10
            }
          ]
        }
      ]
      scale: {
        minReplicas: 0
        maxReplicas: 1
      }
    }
  }
}

output fqdn string = containerApp.properties.configuration.ingress.fqdn
