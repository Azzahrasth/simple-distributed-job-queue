# 🚀 Simple Distributed Job Queue Simulation

A simple asynchronous job queue system built with **Golang**, leveraging **Goroutines** and **Channels** for concurrent job processing and **GraphQL** for the API interface. This project focuses on demonstrating concurrency safety, retry logic, and clean architecture principles.

---

## ✨ Features Implemented

* **Asynchronous Processing:** Jobs are processed in the background by a dedicated worker pool.
* **Concurrency Safety:** Uses `sync.RWMutex` (in Repository) and Go Channels (in Service) to safely handle multiple simultaneous job creations and status updates.
* **Retry Logic:** Failed jobs (like `unstable-job`) automatically retry up to 3 times with **exponential backoff** before being marked as `FAILED`.
* **In-Memory Storage:** Jobs are persisted in a concurrent-safe Go map (in-memory database).
* **GraphQL API:** Provides clear query and mutation endpoints for interaction.

---

## 🛠️ Requirements

* Go >= 1.20
* A running **Redis** instance (optional, for real-world scaling, but not used in this in-memory simulation).

---

## 🏃 How to Run

1.  **Download dependencies and tidy the project:**
    ```bash
    go mod tidy
    ```

2.  **Run the application:**
    ```bash
    go run main.go
    ```

3.  **Access the GraphQL Playground:**
    Open your browser and navigate to: **`http://localhost:58579/graphiql`**

---

## 🌐 GraphQL Endpoints

You can test the core functionalities using the following operations in the GraphiQL explorer:

### Mutations (Creating and Enqueuing Jobs)

| Mutation | Description | Verification |
| :--- | :--- | :--- |
| `SimultaneousCreateJob` | Creates multiple jobs concurrently to test safety. | Status should eventually be `COMPLETED` for all. |
| `SimulateUnstableJob` | Creates a job that fails twice but **succeeds on the 3rd attempt** (Retry Logic Test). | Final status must be `COMPLETED` with `attempts: 3`. |
| `Enqueue` | Creates a single standard job. | |

### Queries (Checking Status)

| Query | Description |
| :--- | :--- |
| `GetAllJobs` | Lists all jobs created in the system with their IDs. |
| `GetJobById` | Fetches the full status and attempt count for a specific job ID. |
| `GetAllJobStatus` | Returns the aggregated count of jobs by status (`PENDING`, `RUNNING`, `FAILED`, `COMPLETED`). |

