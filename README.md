# 🚀 **ProjectTextAgent**

Part of a [`ProjectText Suite`](https://github.com/Flagro/ProjectTextSuite). ProjectTextAgent is an observer for your project directory for updating the Text and Table databases for future question-answering pipelines.

## 🌟 **Features**
- Implements concurrent Goroutines watchers to send updates when files are modified/deleted/created
- On events extracts text/tables from files with TextTableSpoon
- Updates the data in Vector database (VecMetaQ) and simple relational database (PostgreSQL) with new texts

## 🚀 **Getting Started**
Make sure to have docker installed on your system and then simply copy and initialize the .env file and do a docker compose up:
```bash
mv .env-example .env
docker compose up
```

## 📘 **Usage**
Running the image makes VecMetaQ and PostgreSQL databases available at the hosts you specified in the .env file and updates these databases with text contents of files in directory you specified to enable Question-Answering pipelines over these texts with LLMs (check out https://github.com/Flagro/ProjectTextQnA).

## 🤝 **Collaboration & Issues**
Open for collaboration; check the [issues page](https://github.com/Flagro/ProjectTextAgent/issues) for discussions.

## 🌟 **Shoutout to the dependencies**
This project uses:
- https://github.com/fsnotify/fsnotify for file update events handling
- https://gorm.io - orm to connect ProjectTextAgent with postgres for for the flat tables manipulations
- https://github.com/Flagro/VecMetaQ - wrapper over vector database for easy sever-like interface for text embeddings
- https://github.com/Flagro/TextTableScoop - CLI file-to-text parser
