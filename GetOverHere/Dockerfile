FROM maven:3.8.4-openjdk-11-slim AS build
COPY . /home/app/
RUN mvn -f /home/app/pom.xml clean install

FROM openjdk:11-jre-slim
COPY --from=build /home/app/target/libs /app/libs
COPY --from=build /home/app/target/*.jar /app/app.jar
WORKDIR /app
EXPOSE 4567
ENTRYPOINT ["java","-jar","app.jar"]

