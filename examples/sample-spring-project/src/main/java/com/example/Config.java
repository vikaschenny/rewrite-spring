package com.example;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Configuration;

@Configuration
public class Config {

    @Value("${server.port:8080}")
    private int serverPort;

    @Value("${spring.datasource.url}")
    private String databaseUrl;

    @Value("${spring.datasource.username:defaultUser}")
    private String databaseUsername;

    // Old deprecated property that will be transformed
    @Value("${management.metrics.binders.jvm.enabled:true}")
    private boolean jvmMetricsEnabled;

    @Value("${management.metrics.binders.tomcat.enabled:false}")
    private boolean tomcatMetricsEnabled;

    @ConditionalOnProperty(name = "management.metrics.binders.jvm.enabled", havingValue = "true")
    public void configureJvmMetrics() {
        // Configuration logic here
    }

    @ConditionalOnProperty(
        name = "management.endpoints.web.exposure.include",
        havingValue = "health,info,metrics"
    )
    public void configureActuatorEndpoints() {
        // Configuration logic here
    }

    // Getters
    public int getServerPort() {
        return serverPort;
    }

    public String getDatabaseUrl() {
        return databaseUrl;
    }

    public String getDatabaseUsername() {
        return databaseUsername;
    }

    public boolean isJvmMetricsEnabled() {
        return jvmMetricsEnabled;
    }

    public boolean isTomcatMetricsEnabled() {
        return tomcatMetricsEnabled;
    }
} 