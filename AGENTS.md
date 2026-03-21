# AGENTS.md

## Cursor Cloud specific instructions

### Project Overview

Personal blog system (个人博客系统) with two sub-projects:
- **`blog-back/`** — Spring Boot 2.6.13 REST API (Java, Maven, port 9090)
- **`blog-front/`** — Vue 3 + Vite SPA (TypeScript, npm, port 5173)

### Running Services

**MySQL** must be running before the backend starts:
```
sudo service mysql start
```

**Backend** (from `blog-back/`):
```
mvn spring-boot:run -DskipTests
```

**Frontend** (from `blog-front/`):
```
npm run dev
```

### Key Gotchas

- The MySQL JDBC URL requires `&allowPublicKeyRetrieval=true` for local MySQL 8.0 connections. This is already configured in `application.properties`.
- The backend `application.properties` and frontend `src/utils/request.ts` + `vite.config.ts` are configured for `localhost`. The original codebase pointed to a remote server (`47.114.122.96`).
- `npm run build` fails due to a **pre-existing** missing asset (`../assets/github.png` referenced in `src/components/Mine.vue`). This does not affect the dev server.
- TypeScript checking (`vue-tsc -b --noEmit`) passes cleanly.
- The backend compiles and runs on JDK 21 despite `pom.xml` specifying Java 1.8 source/target.
- The upload path is configured to `/tmp/blog-upload/` for local development (original: `/www/wwwroot/rain/upload/`).

### Database Setup

The backend requires a MySQL database named `blog` with tables: `user`, `article`, `tag`, `article_tag`, `hollow`. To recreate:
```sql
CREATE DATABASE IF NOT EXISTS blog CHARACTER SET utf8mb4;
CREATE USER IF NOT EXISTS 'blog'@'localhost' IDENTIFIED BY 'ab2681219';
GRANT ALL PRIVILEGES ON blog.* TO 'blog'@'localhost';
```
Table schemas map directly from entity classes in `blog-back/src/main/java/com/example/blog/entity/`.

### Testing

- **Backend tests**: `cd blog-back && mvn test`
- **Frontend type check**: `cd blog-front && npx vue-tsc -b --noEmit`
- No ESLint config exists in the frontend project.
