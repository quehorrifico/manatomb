# Mana Tomb

Mana Tomb is a full-stack web application for Magic: The Gathering players to build, manage, and share their Commander decks. It features a modern tech stack with a Go backend and a React frontend, designed to be both powerful and scalable.

---

## Features

* **Secure User Authentication**: Full support for user account creation, login, and session management.
* **Card Search**: A powerful and fast interface to search for any MTG card using the Scryfall API.
* **Deck Management**: Users can create, view, edit, and delete their decks.
* **Intuitive Deckbuilding**: A seamless interface on the deck detail page for adding cards to a main deck or a maybeboard.
* **Deck Analysis**: Automatic mana curve and color distribution charts to help users analyze their builds.
* **Public Profiles & Sharing**: Users can make their decks public and share them via a personal profile page.

---

## Tech Stack

| Category      | Technology                               |
| ------------- | ---------------------------------------- |
| **Frontend**  | React, React Router, Axios, Recharts     |
| **Backend**   | Go, Gin Gonic                            |
| **Database**  | PostgreSQL                               |
| **API**       | RESTful API                              |
| **Styling**   | CSS Modules                              |
| **Dev Tools** | Docker, `golang-migrate`                 |

---

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* [Go](https://golang.org/doc/install) (version 1.22 or later)
* [Node.js](https://nodejs.org/) (version 18 or later)
* [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose
* [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI

### Installation

1.  **Clone the repository:**
    ```
    git clone [https://github.com/quehorrifico/manatomb.git](https://github.com/quehorrifico/manatomb.git)
    cd manatomb
    ```

2.  **Backend Setup:**
    * Navigate to the backend directory: `cd backend`
    * Create your environment file: `cp .env.example .env`
    * Install dependencies: `go mod tidy`

3.  **Frontend Setup:**
    * Navigate to the frontend directory: `cd ../frontend`
    * Install dependencies: `npm install`

### Running the Application

You will need three separate terminal windows to run the full application.

1.  **Start the Database:**
    * From the **root** `manatomb` directory, run:
        ```
        docker-compose up -d
        ```

2.  **Run Database Migrations:**
    * Navigate to the `backend` directory.
    * Apply all database migrations:
        ```
        make migrate-up
        ```

3.  **Start the Backend Server:**
    * In the `backend` directory, run:
        ```
        go run main.go
        ```
    * The backend will be available at `http://localhost:8080`.

4.  **Start the Frontend Server:**
    * In the `frontend` directory, run:
        ```
        npm start
        ```
    * The frontend will open in your browser at `http://localhost:3000`.

---

## API Endpoints

A brief overview of the available API endpoints. All `/api/decks` and `/api/users/me` routes require authentication.

| Method   | Endpoint                          | Description                               |
| -------- | --------------------------------- | ----------------------------------------- |
| `POST`   | `/api/users/register`             | Register a new user.                      |
| `POST`   | `/api/users/login`                | Log in a user and create a session.       |
| `POST`   | `/api/users/logout`               | Log out a user and destroy the session.   |
| `GET`    | `/api/users/me`                   | Get the current logged-in user's details. |
| `GET`    | `/api/profiles/:username`         | Get a user's public profile and decks.    |
| `GET`    | `/api/decks`                      | Get all decks for the logged-in user.     |
| `POST`   | `/api/decks`                      | Create a new deck.                        |
| `GET`    | `/api/decks/:deckId`              | Get details for a single deck.            |
| `PUT`    | `/api/decks/:deckId`              | Update a deck's name/description.         |
| `DELETE` | `/api/decks/:deckId`              | Delete a deck.                            |
| `PUT`    | `/api/decks/:deckId/visibility`   | Set a deck's public/private status.       |
| `POST`   | `/api/decks/:deckId/cards`        | Add a card to a deck.                     |
| `DELETE` | `/api/decks/:deckId/cards/:cardId`| Remove a card from a deck.                |

