# TextLib RL System Makefile

.PHONY: build run clean test help logs models monitor health-check

# Default target
help:
	@echo "TextLib RL System - Available targets:"
	@echo "  build         - Build the Docker image"
	@echo "  run           - Run training with default settings"
	@echo "  run-quick     - Run training with reduced episodes for testing"
	@echo "  run-full      - Run full training with all episodes"
	@echo "  clean         - Clean up containers, images, and generated files"
	@echo "  test          - Run basic functionality tests"
	@echo "  health-check  - Perform system health check"
	@echo "  logs          - View training logs"
	@echo "  models        - List available models"
	@echo "  monitor       - Start monitoring dashboard"
	@echo "  report        - Generate API usage report"
	@echo "  setup         - Initial setup and directory creation"

# Build the Docker image
build:
	@echo "Building TextLib RL learner..."
	docker build -t textlib-rl-learner .

# Setup directories and permissions
setup:
	@echo "Setting up directories..."
	mkdir -p logs data models configs
	chmod -R 755 logs data models configs

# Run training with default settings
run: setup build
	@echo "Running RL training with default settings..."
	./run-rl-training.sh

# Quick test run with fewer episodes
run-quick: setup build
	@echo "Running quick training test..."
	./run-rl-training.sh 100 info false

# Full training run
run-full: setup build
	@echo "Running full training..."
	./run-rl-training.sh 10000 info true

# Health check
health-check: build
	@echo "Performing health check..."
	docker run --rm \
		--memory="512m" \
		--cpus="1" \
		textlib-rl-learner \
		./rl-textlib-learner --mode=health-check

# View logs
logs:
	@echo "Available log files:"
	@ls -la logs/ 2>/dev/null || echo "No logs directory found"
	@echo ""
	@echo "Recent log entries:"
	@find logs -name "*.json" -type f -exec echo "=== {} ===" \; -exec head -20 {} \; 2>/dev/null || echo "No log files found"

# List models
models:
	@echo "Available models:"
	@ls -la models/ 2>/dev/null || echo "No models directory found"

# Generate report from existing insights
report:
	@if [ -f "insights.json" ]; then \
		echo "Generating report from existing insights..."; \
		docker run --rm \
			-v "$(PWD)/insights.json:/data/insights.json:ro" \
			-v "$(PWD):/output:rw" \
			textlib-rl-learner \
			./rl-textlib-learner \
			--mode=generate-report \
			--input=/data/insights.json \
			--output=/output/api-usage-guide.md; \
		echo "Report generated: api-usage-guide.md"; \
	else \
		echo "No insights.json found. Run training first."; \
	fi

# Start monitoring with docker-compose
monitor:
	@echo "Starting monitoring stack..."
	docker-compose up -d metrics-collector
	@echo "Prometheus available at: http://localhost:9090"

# Stop monitoring
monitor-stop:
	@echo "Stopping monitoring stack..."
	docker-compose down

# Test the system
test: build health-check
	@echo "Running basic functionality tests..."
	@echo "Test 1: Health check - PASSED"
	
	@echo "Test 2: Quick training run..."
	@./run-rl-training.sh 10 debug false > test_output.log 2>&1 && echo "Test 2: Quick training - PASSED" || echo "Test 2: Quick training - FAILED"
	
	@echo "Test 3: Report generation..."
	@if [ -f "insights.json" ]; then \
		make report > /dev/null 2>&1 && echo "Test 3: Report generation - PASSED" || echo "Test 3: Report generation - FAILED"; \
	else \
		echo "Test 3: Report generation - SKIPPED (no insights)"; \
	fi
	
	@rm -f test_output.log

# Clean up everything
clean:
	@echo "Cleaning up..."
	@docker stop textlib-rl-training 2>/dev/null || true
	@docker rm textlib-rl-training 2>/dev/null || true
	@docker-compose down 2>/dev/null || true
	@docker rmi textlib-rl-learner 2>/dev/null || true
	@echo "Cleaning generated files..."
	@rm -rf logs/* data/* models/* 2>/dev/null || true
	@rm -f insights.json final_model.json api-usage-guide.md 2>/dev/null || true
	@echo "Cleanup complete"

# Deep clean including Docker volumes
clean-all: clean
	@echo "Performing deep clean..."
	@docker system prune -f
	@docker volume prune -f

# Show system status
status:
	@echo "=== System Status ==="
	@echo "Docker:"
	@docker --version 2>/dev/null || echo "  Docker: Not available"
	@echo ""
	@echo "Containers:"
	@docker ps -a --filter name=textlib 2>/dev/null || echo "  No textlib containers"
	@echo ""
	@echo "Images:"
	@docker images --filter reference=textlib-rl-learner 2>/dev/null || echo "  No textlib images"
	@echo ""
	@echo "Directories:"
	@ls -la logs data models 2>/dev/null || echo "  Directories not created yet"
	@echo ""
	@echo "Generated Files:"
	@ls -la *.json *.md 2>/dev/null || echo "  No generated files"

# Development mode - run with live code mounting
dev: setup build
	@echo "Running in development mode..."
	docker run -it --rm \
		--name textlib-rl-dev \
		--memory="512m" \
		--cpus="2.0" \
		-v "$(PWD)/logs:/app/logs:rw" \
		-v "$(PWD)/data:/app/data:rw" \
		-v "$(PWD)/models:/app/models:rw" \
		-e LOG_LEVEL=debug \
		-e MAX_EPISODES=50 \
		textlib-rl-learner \
		/bin/sh

# Quick development test
dev-test: setup build
	@echo "Running development test..."
	docker run --rm \
		--memory="256m" \
		--cpus="1.0" \
		-v "$(PWD)/logs:/app/logs:rw" \
		-v "$(PWD)/data:/app/data:rw" \
		-v "$(PWD)/models:/app/models:rw" \
		-e LOG_LEVEL=debug \
		textlib-rl-learner \
		./rl-textlib-learner --mode=train --episodes=5