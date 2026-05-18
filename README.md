# example-wikipedia-scraper

Example Wikipedia Scraper is a simple web scraper built using Go lang and Docker. It is designed to scrape data from Polish Wikipedia pages and store it in a structured format for further analysis.
It can be expanded to scrape data from other webpages with listed supage links. The scraper is designed to be easy to use and can be run using Docker Compose, making it easy to deploy and manage.

Technologies Used:
- Go lang: The scraper is built using Go lang, which is a powerful and efficient programming language that is well-suited for web scraping tasks.
- Docker: The scraper is containerized using Docker, which allows for easy deployment and management of the application.
- Docker Compose: The scraper can be run using Docker Compose, which allows for easy orchestration of the different components of the application.
- Chrome DevTools Protocol: The scraper uses the Chrome DevTools Protocol to interact with the web pages and extract data.
- PostgreSQL: The scraper stores the scraped data in a PostgreSQL database for further analysis.
- RabbitMQ: The scraper uses RabbitMQ for message queuing and task management.

Features:
- Scrapes data from Polish Wikipedia pages and stores it in a structured format.
- Can be easily expanded to scrape data from other webpages with listed subpage links.
- Containerized using Docker for easy deployment and management.
- Uses Chrome DevTools Protocol for efficient web scraping.
- Stores scraped data in a PostgreSQL database for further analysis.
- Uses RabbitMQ for message queuing and task management.


Requirements:
- Docker: The scraper requires Docker to be installed on the host machine in order to run the application.
- Docker Compose: The scraper requires Docker Compose to be installed on the host machine in order to orchestrate the different components of the application.
- Available port 8080: The scraper uses port 8080 to serve the application, so it must be available on the host machine.

Usage:
1. Clone the repository:
2. Navigate to the project directory
3. Build and run the application using `bin/run` or `bin/run-local` scripts.
4. The scraper will start and begin scraping data from the specified Wikipedia pages.
5. The scraped data will be stored in the PostgreSQL database for further analysis.
6. Go to `http://localhost:8080/adminer/` to access the adminer with database and view the scraped data. Creditials: `user` `password` host `db`