# Stratal🏗️

**Stratal** 🏗️ is a robust, distributed job orchestration and workflow management system built in Go with a modern React web interface. It's designed to handle complex task execution workflows with support for parallel processing, task dependencies, and secure secret management.

## Core Architecture

**Stratal** has three main components:

1. **HTTP Server** (`cmd/server`) - REST API for job management and web UI
2. **Worker Processes** (`cmd/worker`) - Distributed task execution engines  
3. **CLI Tool** (`cmd/cli`) - Command-line interface for system interaction

## Key Features

### 🔄 **Advanced Job Processing**
- **Dependency-aware task execution** with topological sorting
- **Parallel task processing** - tasks with the same order level run concurrently
- **Task output passing** - outputs from completed tasks are available as environment variables to subsequent tasks
- **Multiple execution engines** supporting both custom scripts and built-in tasks

### 🛡️ **Enterprise Security**
- **AES encryption** for sensitive data and secrets
- **Secure parameter injection** into task environments
- **User-based secret isolation** with encrypted storage

### 📊 **Comprehensive Monitoring**
- **Real-time logging system** with structured job run logs
- **WebSocket-based log streaming** for live monitoring
- **Job status tracking** (pending → queued → running → completed/failed)
- **Task-level execution tracking** with detailed error reporting

### 🔧 **Flexible Task Types**

**Custom Scripts**: Execute code in multiple languages
```json
{
  "type": "custom",
  "config": {
    "script": {
      "language": "python|javascript|bash|go|ruby|php|perl",
      "code": "your script here"
    }
  }
}
```

**Built-in Tasks**: Pre-built integrations
- `http_request` - REST API calls with full HTTP method support
- `send_email` - SMTP email delivery
- `format_output` - Data transformation and formatting
- `echo` - Simple testing and debugging

### 🏗️ **Production-Ready Infrastructure**
- **Redis-based job queue** for reliable task distribution
- **PostgreSQL storage** with comprehensive database schema
- **Database migrations** using golang-migrate
- **Graceful shutdown** handling for all components
- **Connection pooling** with pgx for optimal database performance

### 🌐 **Modern Web Interface**
- **React + TypeScript** frontend with Vite build system
- **Real-time dashboard** showing job statistics and system health
- **Job creation wizard** with visual task configuration
- **Live log viewer** with WebSocket streaming
- **Responsive design** for desktop and mobile access

### 🔌 **Developer Experience**
- **Comprehensive REST API** with detailed examples
- **Insomnia test collection** for API testing
- **Docker support** for development and deployment
- **Makefile** with common development tasks
- **SQLC integration** for type-safe database operations

## Use Cases

**Stratal** excels at:
- **Data pipeline orchestration** with complex dependencies
- **Automated deployment workflows** with rollback capabilities  
- **ETL/ELT processing** with multi-stage transformations
- **Scheduled batch processing** with cron-like scheduling
- **Microservice coordination** in distributed systems
- **CI/CD pipeline execution** with parallel build stages

The system is particularly well-suited for organizations needing reliable, scalable job processing with strong security requirements and comprehensive monitoring capabilities.

## Future Work

To transform Stratal into a comprehensive automation platform, the following enhancements are planned:

### 🕐 **Scheduling & Triggers**
- **Cron-based scheduling**: Add cron expressions to jobs for automatic execution
- **Event-driven triggers**: Webhook endpoints, file system watchers, database changes
- **Job chaining**: Trigger jobs based on completion of other jobs
- **Retry mechanisms**: Configurable retry policies with exponential backoff

### ⚡ **Enhanced Task Types**
- **File operations**: Copy, move, delete, compress, archive files
- **Database tasks**: SQL execution, data migration, backup/restore
- **Cloud integrations**: AWS S3, Azure Blob, GCP operations
- **Git operations**: Clone, pull, push, tag repositories
- **Docker/Kubernetes**: Container management and deployment tasks
- **Notification tasks**: Slack, Teams, Discord, SMS alerts

### 🔀 **Workflow Improvements**
- **Conditional execution**: If/else logic based on task outputs or environment
- **Loop constructs**: Iterate over arrays/lists in workflows
- **Template variables**: Dynamic parameter substitution across jobs
- **Job templates**: Reusable job definitions with parameters
- **Nested workflows**: Sub-jobs and job composition

### 📈 **Monitoring & Observability**
- **Metrics dashboard**: Job success rates, execution times, resource usage
- **Alerting system**: Configurable alerts for failures, long-running jobs
- **Audit logging**: Comprehensive job execution history
- **Performance profiling**: Resource consumption tracking
- **Health checks**: System and service health monitoring

### 💻 **CLI Enhancements**
- **Job management**: Create, update, delete, and monitor jobs via CLI
- **Job execution**: Run jobs and view real-time status
- **Configuration management**: Manage system settings and environments
- **System health**: Check service status and diagnostics
- **Import/export**: Backup and restore job definitions

### ⚙️ **Configuration & Management**
- **YAML/JSON job definitions**: File-based job configuration
- **Environment management**: Dev/staging/prod configurations  
- **Secret rotation**: Automatic credential updates
- **Backup/restore**: System state management
- **Multi-tenancy**: User isolation and resource quotas

### 🔗 **Integration Capabilities**
- **API rate limiting**: Built-in throttling for external services
- **Circuit breakers**: Fault tolerance for external dependencies
- **Message queues**: RabbitMQ, Kafka integration
- **Monitoring tools**: Prometheus, Grafana integration
- **CI/CD pipelines**: Jenkins, GitHub Actions integration

### 🎯 **High Priority Quick Wins**
1. **Job scheduling** - Leverage existing cron dependency for automated execution
2. **Enhanced CLI** - Build comprehensive command interface
3. **File operations** - Add basic file manipulation tasks
4. **Conditional logic** - Add if/else capabilities to workflows
5. **Job templates** - Create reusable job patterns

These improvements will transform Stratal from a job orchestrator into a full-featured automation platform suitable for DevOps, data engineering, and general workflow automation needs.

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, improving documentation, or suggesting enhancements, your help is appreciated.

Please read our [Contributing Guidelines](Contributing.md) for details on:
- Setting up your development environment
- Code style and standards
- Submitting pull requests
- Reporting bugs and requesting features

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
