# Build the docker image
FROM maven:3.9-eclipse-temurin-21-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the project files
COPY pom.xml .
COPY src ./src
COPY config ./config

# Set property files
RUN mv src/main/resources/prod-application.example.yml src/main/resources/application.yml

# Run the maven lifecycles to build the application
RUN mvn clean package -DskipTests

# Stage 2: Create the run time image
FROM eclipse-temurin:21-jdk-alpine AS run

# Set the working directory
WORKDIR /app

RUN apk add --no-cache curl jq

HEALTHCHECK CMD curl -sSf http://127.0.0.1:8080/details/health | jq '.status' | grep -q '\"UP\"' || exit 1


# Copy the JAR file built in the previous stage
COPY --from=build /app/target/*.jar ./app.jar

ENV PORT=8080

# Expose the applications port
EXPOSE ${PORT}

# Command to run the application
CMD ["java", "-jar", "app.jar"]