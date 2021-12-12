package com.zubankov.GetOverHere;

import com.google.gson.Gson;
import javax.servlet.MultipartConfigElement;
import javax.servlet.http.Part;
import java.io.File;
import java.io.InputStream;
import java.nio.file.FileAlreadyExistsException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import static spark.Spark.*;

public class Main {

    public static final String FILE_ALREADY_EXISTS = "File Already exists at storage device";
    public static final String FILE_UPLOADED = "File uploaded successfully";

    public static void main(String[] args) {
        post("/upload", "multipart/form-data", (request, response) -> {
            Gson responseJson = new Gson();
            String location = "files";          // the directory location where files will be stored
            long maxFileSize = 100000000;       // the maximum size allowed for uploaded files
            long maxRequestSize = 100000000;    // the maximum size allowed for multipart/form-data requests
            int fileSizeThreshold = 1024;       // the size threshold after which files will be written to disk

            MultipartConfigElement multipartConfigElement = new MultipartConfigElement(
                    location, maxFileSize, maxRequestSize, fileSizeThreshold);
            request.raw().setAttribute("org.eclipse.jetty.multipartConfig",
                    multipartConfigElement);

            String fName = request.raw().getPart("file").getSubmittedFileName();
            System.out.println("File: " + fName);

            Part uploadedFile = request.raw().getPart("file");
            File directory = new File(location);
            if (! directory.exists()){
                directory.mkdirs();
            }
            Path out = Paths.get(location + "/" + fName);
            try (final InputStream in = uploadedFile.getInputStream()) {
                Files.copy(in, out);
                uploadedFile.delete();
            }
            catch (FileAlreadyExistsException e){
                response.status(400);
                System.err.println(responseJson.toJson(new Log("error", FILE_ALREADY_EXISTS)));
                return responseJson.toJson(new Message(FILE_ALREADY_EXISTS));
            }
            System.out.println(responseJson.toJson(new Log("info", FILE_UPLOADED)));
            return responseJson.toJson(new Message(FILE_UPLOADED));
        });
        after(((request, response) -> {
            response.type("application/json");
        }));
    }
}
