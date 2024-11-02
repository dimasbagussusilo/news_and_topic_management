News and Topic Management API
=============================

This project implements a REST API for managing news articles and topics, allowing you to perform CRUD operations and manage relationships between news and topics.

Table of Contents
-----------------

*   [Getting Started](#getting-started)
*   [Prerequisites](#prerequisites)
*   [Installation](#installation)
*   [Common Tasks](#common-tasks)
*   [API Documentation](#api-documentation)
*   [Endpoints](#endpoints)
*   [Testing](#testing)
*   [License](#license)

Getting Started
---------------

### Prerequisites

*   **Go** (version 1.20)
*   **Docker**
*   **Makefile**
*   **Go AIR**
*   [ZOG Moonrepo](https://github.com/zero-one-group/monorepo)

### Installation

1.  **Clone the repository:**

        git clone https://github.com/dimasbagussusilo/news_and_topic_management
        cd news_and_topic_management

2.  **Make sure Docker is running and you have Moonrepo installed.**

Common Tasks
------------

### Building the Project

To build and run the necessary components (like PostgreSQL and Go libraries), use:

    moon :up

### Stopping the Project

To stop all running services:

    moon :down

### Destroying the Project

To destroy all resources created during `up`:

    moon :destroy

### Running Tests

To run tests for everything in the project:

    moon :tests

API Documentation
-----------------

Explore our interactive Swagger documentation to unlock the full potential of our API!

Visit:

    {{BASE_URL}}/swagger/index.html

Just replace `{{BASE_URL}}` with your base URL and start exploring!

Endpoints
---------

### News Endpoints

*   **GET /news**
    *   Retrieve all news articles.
    *   **Query Parameters:**
        *   `title` (optional): Filter by Title that contain the input.
        *   `status` (optional): Filter by status (draft, deleted, published).
        *   `start_date` (optional): Filter by range date.
        *   `end_date` (optional): Filter by range date.
        *   `limit` (optional): Limit data that you need.
        *   `page` (optional): Set current page data.
        *   `sort_by` (optional): Sort data by input.
        *   `sort_order` (optional): Sort data by 'asc' or 'desc'.


*   **POST /news**
    *   Create a new news article.
      *   **Request Body:**

              {
                  "title": "Covid 19 is gone!",
                  "content": "Alhamdulillah the covid 19 pandemic is over, we can continue our activities without fear of catching a virus",
                  "author_id": 1,
                  "status": "test",
                  "topic_ids": [
                       1,
                       5
                 ]
              }

*   **GET /news/{id}**
    *   Retrieve a specific news article by ID.
*   **PUT /news/{id}**
    *   Update an existing news article.
    *   **Request Body:**

            {
                 "title": "Updated Title",
                 "content": "Updated Content",
                 "status": "published"
            }

*   **DELETE /news/{id}**
    *   Delete a specific news article by ID.

### Topic Endpoints

*   **GET /topics**
    *   Retrieve all topics.
*   **POST /topics**
    *   Create a new topic.
    *   **Request Body:**

            {
                "name": "Beauty"
            }

*   **GET /topics/{id}**
    *   Retrieve a specific topic by ID.
*   **PUT /topics/{id}**
    *   Update an existing topic.
    *   **Request Body:**

            {
                "name": "Updated Topic Name"
            }

*   **DELETE /topics/{id}**
    *   Delete a specific topic by ID.

Testing
-------

Your project includes a test suite to ensure functionality. To run the tests, use:

    moon :tests

Make sure that your testing framework is correctly set up in the project. You can add new tests by creating files in the `tests` directory.

Diagram Project
-------
![project diagram](https://raw.githubusercontent.com/dimasbagussusilo/news_and_topic_management/refs/heads/main/mermaid-diagran.png)
![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

Database Schema
-------
![db diagram](https://raw.githubusercontent.com/dimasbagussusilo/news_and_topic_management/refs/heads/main/db-diagram.png)

License
-------

This project is licensed under the MIT License.
