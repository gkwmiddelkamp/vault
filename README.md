<img src="https://vault.pnck.nl/cdn/logo.png"  alt=""/>

Vault is a project to have a light-weight, secure and multi-tenant solution for encrypted password storage. It provides a simple Rest API where you can manage your environments and tokens.
It uses a MongoDB database as the storage backend. 

**This project is a work in progress, do not use in production until v1.0.0 is available**

# Environments
Security is key in the project. You can create separate environments for your projects or customers. All environments use unique encryption keys, which are never stored in the database and are only available to the customer.
At the first start of the application, the Master Admin token will be logged as output once. Save it, it will never be shown again.

If you missed the token after the first start, you need to remove the collections (environment, token, secret) from the database and restart the application. None of the tokens are recoverable.

The MasterAdmin token can create an Environment. As a response to this call an EnvironmentAdmin token is returned once. This type of token can be used to create ReadWrite or ReadOnly tokens. Read the section [Tokens](#Tokens) for more detailed view of the different token types.

# Tokens
There are 4 types of tokens, each having its own purpose.


|                                  | MasterAdmin   	 | EnvironmentAdmin  	 | ReadWrite  	 | ReadOnly   	 |
|----------------------------------|-----------------|---------------------|--------------|--------------|
| Create MasterAdmin token	        | 	     ✅         | 	                   | 	            | 	            |
| Create EnvironmentAdmin token	   | 	 ✅              | 	                   | 	            | 	            |
| Create ReadWrite/ReadOnly token	 | 	               | 	      ✅            | 	            | 	            |
| Manage environments              | 	     ✅          | 	                   | 	            | 	            |
| Manage secrets	                  | 	               | 	                   | 	  ✅          | 	            |
| Get decrypted secret             | 	               | 	                   | 	   ✅         | 	    ✅        |


# Getting started
Vault can be run as a stand-alone application on a server, or run as a container in Docker or Kubernetes.

Make sure you always run the latest release version. 

The entire application is built stateless and supports multiple replicas for load balancing and high-availability purposes.

## Networking
Vault itself does not handle TLS traffic. The service that exposes the application has to handle and forward to the application port.

## Configuration
The application can be configured using environment variables for the database connection.

| Parameter               | Description   	                                                                                                                                                          | Default  	 |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------|
| ```PORT```	             | Port for the REST API endpoints	                                                                                                                                         | 	 ```8080``` |
| ```MONGODB_URI```	      | MongoDB connect URI to connect to the database<br/>Example: ```mongodb+srv://username:password@database-host/database-name?retryWrites=true&w=majority&appName=Vault```	 | 	     |
| ```MONGODB_DATABASE```	 | Database name if not provided in the connect URI<br/>Example: ```vault```                                                                                                   | 	  |

## Kubernetes deployment
```yaml
apiVersion: v1
data:
  MONGODB_URI: bW9uZ29kYitzcnY6Ly91c2VybmFtZTpwYXNzd29yZEBkYXRhYmFzZS1ob3N0L2RhdGFiYXNlLW5hbWU/cmV0cnlXcml0ZXM9dHJ1ZSZ3PW1ham9yaXR5JmFwcE5hbWU9VmF1bHQ=
kind: Secret
metadata:
  name: vault-database
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: vault
  name: vault
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vault
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: vault
    spec:
      containers:
      - image: ghcr.io/gkwmiddelkamp/vault
        name: vault
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: vault-database
              key: MONGODB_URI
        resources:
            requests:
              cpu: 50m
              memory: 200Mi
        securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            runAsUser: 1000
            seccompProfile:
              type: RuntimeDefault

```