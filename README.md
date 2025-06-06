# Notifier - Payment Notifier
Notifier is an application written in Go designed to automatically notify payments using cron jobs. This project allows you to schedule recurring tasks that send notifications about pending payments, payments made or any event related to payments.

## Characteristics
- Automation of payment notifications through cron jobs.
- Flexible configuration for different execution frequencies.
- Integration with notification email services.
- Detailed logs for monitoring and debugging.

## Requirements
- Go 1.20 or higher.
- Access to a supported MySQL database.
- Configuration of the cron service to execute recurring tasks.

## Installation
1. Clone this repository:
    ```bash
    git clone https://github.com/rober0xf/notifier.git
    cd notifier
    ```
2. Init the project:
  ```bash
  go mod init github.com/rober0xf/notifier
  ```
3. Install the dependencyes:
    ```bash
    go mod tidy
    ```
4. Configure the env variables in a file `.env`:
    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=notifier
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=mail@example.com
    SMTP_PASSWORD=your_password
    ```

## Usage
1. Compile the application:
    ```bash
    go build -o notifier
    ```
2. Execute the binary:
    ```bash
    ./notifier
    ```
3. Configure a cronjob to execute the application. For example each hour:
    ```
    0 * * * * /route/to/notifier >> /var/log/notifier.log 2>&1
    ```

---
