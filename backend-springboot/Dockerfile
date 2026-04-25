#base java 17 image
FROM eclipse-temurin:17-jdk

#sets the working folder inside the container
WORKDIR /app

#copues maven wrapper, project config, and source code into the image
COPY .mvn/ .mvn
COPY mvnw pom.xml ./
COPY src ./src

#makes the Maven wrapper executable
RUN chmod +x mvnw
#builds the spring boot app into a jar file and skips the test
RUN ./mvnw clean package -DskipTests

#documents that the app listens to port 8080
EXPOSE 8080

#starts the built jar when the container runs
CMD ["java", "-jar", "target/backend-springboot-0.0.1-SNAPSHOT.jar"]

#Jar - Jar Archive - packaged app file that java can run
#Maven - java build tool, that runs dependencies, builds the project, runs tests, and packages the app, equivalent to npm in node.js
