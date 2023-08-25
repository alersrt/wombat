package org.student.wombat;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.lang.invoke.MethodHandles;

/**
 * The main class of the application.
 */
@SpringBootApplication
public class WombatApplication {

    private static final Logger log = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());

    /**
     * The endpoint of the application.
     *
     * @param args arguments of this app.
     */
    public static void main(String... args) {
        log.info("*** Starting application...");
        SpringApplication.run(WombatApplication.class, args);
        log.info("*** Application started...");
    }
}
