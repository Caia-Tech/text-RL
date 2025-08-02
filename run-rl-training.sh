#!/bin/bash

# TextLib RL Training Execution Script
# This script safely runs the RL training system in an isolated container environment

set -e

echo "Starting TextLib RL Training System..."

# Configuration
IMAGE_NAME="textlib-rl-learner"
CONTAINER_NAME="textlib-rl-training"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed or not in PATH"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    print_error "Docker daemon is not running"
    exit 1
fi

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        print_status "Stopping container..."
        docker stop $CONTAINER_NAME || true
    fi
    
    if docker ps -aq -f name=$CONTAINER_NAME | grep -q .; then
        print_status "Removing container..."
        docker rm $CONTAINER_NAME || true
    fi
}

# Set up cleanup trap
trap cleanup EXIT

# Create necessary directories
print_status "Creating necessary directories..."
mkdir -p logs data models configs

# Set permissions
print_status "Setting directory permissions..."
chmod -R 755 logs data models configs

# Build the Docker image
print_status "Building RL training container..."
docker build -t $IMAGE_NAME . || {
    print_error "Failed to build Docker image"
    exit 1
}

# Parse command line arguments
MAX_EPISODES=${1:-10000}
LOG_LEVEL=${2:-info}
ENABLE_PROFILING=${3:-false}

print_status "Configuration:"
echo "  - Max Episodes: $MAX_EPISODES"
echo "  - Log Level: $LOG_LEVEL"
echo "  - Profiling: $ENABLE_PROFILING"

# Run the container with safety constraints
print_status "Starting RL training in safe container..."
docker run \
    --name $CONTAINER_NAME \
    --memory="512m" \
    --memory-swap="512m" \
    --cpus="2.0" \
    --pids-limit 100 \
    --read-only \
    --security-opt=no-new-privileges \
    --cap-drop=ALL \
    --tmpfs /tmp:rw,noexec,nosuid,size=100m \
    -v "$(pwd)/logs:/app/logs:rw" \
    -v "$(pwd)/data:/app/data:rw" \
    -v "$(pwd)/models:/app/models:rw" \
    -e LOG_LEVEL=$LOG_LEVEL \
    -e ENABLE_PROFILING=$ENABLE_PROFILING \
    -e MAX_EPISODES=$MAX_EPISODES \
    -e CHECKPOINT_INTERVAL=500 \
    -e GOMAXPROCS=2 \
    -e GOMEMLIMIT=512MiB \
    $IMAGE_NAME \
    ./rl-textlib-learner --mode=train --episodes=$MAX_EPISODES --log-level=$LOG_LEVEL || {
    print_error "Training failed"
    exit 1
}

# Extract insights after training
print_status "Extracting insights..."
if docker ps -aq -f name=$CONTAINER_NAME | grep -q .; then
    docker cp $CONTAINER_NAME:/app/logs/insights.json ./insights.json 2>/dev/null || {
        print_warning "Could not extract insights.json (may not exist yet)"
    }
    
    # Find the latest model file
    LATEST_MODEL=$(docker exec $CONTAINER_NAME find /app/models -name "final_model_*.json" -type f -printf '%T@ %p\n' 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)
    if [ ! -z "$LATEST_MODEL" ]; then
        docker cp $CONTAINER_NAME:$LATEST_MODEL ./final_model.json
        print_status "Extracted model: $LATEST_MODEL"
    else
        print_warning "No model files found"
    fi
fi

# Generate API usage report
print_status "Generating API usage recommendations..."
if [ -f "./insights.json" ]; then
    docker run --rm \
        -v "$(pwd)/insights.json:/data/insights.json:ro" \
        -v "$(pwd)/final_model.json:/data/model.json:ro" \
        -v "$(pwd):/output:rw" \
        $IMAGE_NAME \
        ./rl-textlib-learner \
        --mode=generate-report \
        --input=/data/insights.json \
        --model=/data/model.json \
        --output=/output/api-usage-guide.md || {
        print_warning "Failed to generate report, trying without model file..."
        docker run --rm \
            -v "$(pwd)/insights.json:/data/insights.json:ro" \
            -v "$(pwd):/output:rw" \
            $IMAGE_NAME \
            ./rl-textlib-learner \
            --mode=generate-report \
            --input=/data/insights.json \
            --output=/output/api-usage-guide.md
    }
    
    if [ -f "./api-usage-guide.md" ]; then
        print_status "API usage guide generated: api-usage-guide.md"
    else
        print_warning "Failed to generate API usage guide"
    fi
else
    print_warning "No insights file found, skipping report generation"
fi

# Display summary
print_status "Training Summary:"
echo "===================="
if [ -f "./logs" ] && [ "$(ls -A ./logs 2>/dev/null)" ]; then
    echo "Log files: $(ls ./logs | wc -l) files"
else
    echo "Log files: 0 files"
fi

if [ -f "./models" ] && [ "$(ls -A ./models 2>/dev/null)" ]; then
    echo "Model files: $(ls ./models | wc -l) files"
else
    echo "Model files: 0 files"
fi

if [ -f "./insights.json" ]; then
    echo "Insights: ✓ Generated"
else
    echo "Insights: ✗ Not found"
fi

if [ -f "./api-usage-guide.md" ]; then
    echo "Usage Guide: ✓ Generated"
else
    echo "Usage Guide: ✗ Not found"
fi

print_status "Training completed successfully!"
print_status "Check the following files for results:"
echo "  - ./logs/ - Training logs and metrics"
echo "  - ./models/ - Trained models and checkpoints"
echo "  - ./insights.json - Detailed insights data"
echo "  - ./api-usage-guide.md - Human-readable usage guide"