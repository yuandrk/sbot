# Telegram Bot for Writing Messages 

####  documentation in process 


This simple Go-based Telegram bot allows you to capture incoming messages and log them into a CSV file. 
Useful for logging or archiving conversations and messages received on Telegram.

Example: 
Used for write data from Apple Shortcuts 


## Continuous Integration and Continuous Deployment Workflow

This document outlines the CI/CD process configured via GitHub Workflows. The workflow is divided into two main stages: Continuous Integration (CI) and Continuous Deployment (CD).

### Continuous Integration (CI)

The CI process includes the following steps:

1. **Action CheckOut**: 
   - Description: Checks out the source code for the workflow.

2. **Run Tests**: 
   - Description: Executes the predefined tests to ensure code reliability and stability.

3. **Login to GHCR.io (GitHub Container Registry)**: 
   - Description: Authenticates to GitHub Container Registry to enable subsequent image push.

4. **Build and Push Docker Image**:
   - Description: Builds the Docker image and pushes it to GitHub Container Registry.
   - Registry: GHCR.io

5. **Clean**: 
   - Description: Performs cleanup operations post-build and push.

### Continuous Deployment (CD)

The CD process encompasses the following steps:

1. **Save Tag and Revision to \$GITHUB_ENV**:
   - Description: Captures and saves the image tag and revision information to the GitHub environment variables (\$GITHUB_ENV).

2. **Update Image Tag in values.yaml**:
   - Description: Updates the `image.tag` field in the `values.yaml` file with the new image tag.
   - File: `values.yaml`

3. **Commit Changes and Push**:
   - Description: Commits the changes made to the `values.yaml` file and pushes them to the repository.
   - Note: This step finalizes the deployment process.

### Example of GitHub WorkFlow

![workflow](/doc/img/mvp_4.png)

### Example of ArgoCD deployment 

![ArgoCD](/doc/img/mvp_3.png)

## Installation Guide

### Prerequisites

- Go (Golang) installed on your system
- Git for cloning the repository
- Telegram account for creating a bot and obtaining API token

### Installation Steps

1. **Clone this repository** by running the following command in your terminal:

`$ git clone git@github.com:andriuk/sbot.git`

2. **Create your own Telegram bot:**

- Go to [BotFather](https://t.me/BotFather) on Telegram.
- Use the `/newbot` command to create a new bot and follow the instructions to get the API token.

Here you go make file which help you build go program 
doc in processing 


3. **Store your Telegram bot API token in an environment variable:**

- Use the command below to set the API token as an environment variable. Replace `<YOUR_API_TOKEN>` with the obtained token.


### Paste your API token
For more secure use command below:

   ` $ read -s TELE_TOKEN`

Paste your API token 

   `$ export $TELE_TOKEN`

4. **Run the bot** using the following command:

`$ ./sbot start`

5. **Interact with your Telegram bot** by sending messages to log them in the CSV file.


#### 6. Well done 

---

If you have any questions or need assistance, feel free to create an issue in this repository. 
We're here to help!