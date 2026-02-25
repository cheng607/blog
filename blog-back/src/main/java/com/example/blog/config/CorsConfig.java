package com.example.blog.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class CorsConfig implements WebMvcConfigurer {

    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**") // 允许所有接口路径
                .allowedOrigins(
                        "http://47.114.122.96",
                        "https://www.raincode.top",
                        "http://www.raincode.top",
                        "https://raincode.top",
                        "http://raincode.top",
                        "http://localhost:5173") // 允许前端的域名（注意不要加端口，因为前端用的是默认80端口）
                .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS") // 允许的请求方法
                .allowedHeaders("*") // 允许的请求头
                .allowCredentials(true) // 允许携带Cookie
                .maxAge(3600); // 预检请求的有效期（秒）
    }
}