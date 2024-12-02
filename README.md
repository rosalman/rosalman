# Forum Project

A simple web forum built in Go with SQLite as the database, allowing users to communicate through posts and comments. This forum also supports features like categories, likes/dislikes, filtering posts, and user authentication. The project is containerized using Docker.

---

## Table of Contents
1. [Features](#features)
2. [Technologies Used](#technologies-used)
3. [Setup Instructions](#setup-instructions)
4. [Database Schema](#database-schema)
5. [Usage](#usage)
6. [Docker Integration](#docker-integration)
7. [Future Improvements](#future-improvements)

---

## Features
### User Authentication
- **Register**: Users can register with email, username, and password.
  - Duplicate email check.
  - Password encryption (Bonus task).
- **Login**: Users can log in with valid credentials.
- **Sessions**: Each session is managed with cookies and has an expiration time.

### Posts & Comments
- Registered users can create posts and comments.
- Posts can be associated with one or more categories.
- Posts and comments are visible to all users (registered and non-registered).

### Likes & Dislikes
- Registered users can like or dislike posts and comments.
- All users can see the number of likes and dislikes.

### Filtering
- Filter posts by:
  - Categories (subforums).
  - Created posts (only visible to the logged-in user).
  - Liked posts (only visible to the logged-in user).

---

## Technologies Used
- **Backend**: Go (Golang)
- **Database**: SQLite
- **Encryption**: bcrypt (Bonus task)
- **Unique Identifiers**: UUID (Bonus task)
- **Containerization**: Docker
- **Testing**: Unit tests in Go

---

## Setup Instructions

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) installed on your system.
- Basic knowledge of Go and SQLite.

### Steps
1. Clone the repository:
   ```bash
   git clone https://learn.reboot01.com/git/mshaban/forum.git
   cd forum
